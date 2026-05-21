import { toast } from "vue-sonner";

export interface AppNotification {
  id: string;
  title: string;
  body: string;
  type: string;
  read: boolean;
  timestamp: Date;
  route: string;
}

const NOTIFICATION_ROUTES: Record<string, string> = {
  task_sent: "/tasks",
  task_graded: "/submissions",
  assessment_complete: "/applications",
  status_update: "/applications",
  interview_scheduled: "/interviews",
  interview_reminder: "/interviews",
  new_application: "/jobs",
  new_submission: "/submissions",
};

interface WSMessage {
  Type: string;
  Text: string;
  Raw: unknown;
}

const RECONNECT_DELAY = 3000;

// Module-level singleton: only one WebSocket per page load regardless of how
// many times useNotifications() is called (HMR, multiple components, etc.)
let _ws: WebSocket | null = null;
let _reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let _intentionalClose = false;

export function useNotifications() {
  const notifications = useState<AppNotification[]>("ws-notifications", () => []);
  const unreadCount = computed(() =>
    notifications.value.filter((n) => !n.read).length,
  );

  const config = useRuntimeConfig();

  function getWSUrl(): string {
    const base = (config.public.apiBase as string) || "http://localhost:4000";
    return base.replace(/^https/, "wss").replace(/^http/, "ws") + "/api/ws";
  }

  function connect() {
    if (import.meta.server) return;
    if (
      _ws &&
      (_ws.readyState === WebSocket.CONNECTING || _ws.readyState === WebSocket.OPEN)
    )
      return;

    _intentionalClose = false;

    try {
      _ws = new WebSocket(getWSUrl());

      _ws.onopen = () => {
        if (_reconnectTimer) {
          clearTimeout(_reconnectTimer);
          _reconnectTimer = null;
        }
      };

      _ws.onmessage = (event: MessageEvent) => {
        try {
          const msg: WSMessage = JSON.parse(event.data as string);

          if (msg.Type === "PING") {
            _ws?.send(JSON.stringify({ Type: "PONG", Text: "Nothing" }));
            return;
          }

          if (msg.Type === "Connected" || !msg.Text) return;

          if (msg.Type === "Notification") {
            try {
              const payload = JSON.parse(msg.Text) as {
                title: string;
                body: string;
                type: string;
              };
              addNotification(payload);
            } catch {
              /* non-JSON text, ignore */
            }
          }
        } catch {
          /* invalid message, ignore */
        }
      };

      _ws.onclose = () => {
        _ws = null;
        if (!_intentionalClose) scheduleReconnect();
      };

      _ws.onerror = () => {
        _ws?.close();
      };
    } catch {
      scheduleReconnect();
    }
  }

  function scheduleReconnect() {
    if (_reconnectTimer) return;
    _reconnectTimer = setTimeout(() => {
      _reconnectTimer = null;
      connect();
    }, RECONNECT_DELAY);
  }

  function disconnect() {
    _intentionalClose = true;
    if (_reconnectTimer) {
      clearTimeout(_reconnectTimer);
      _reconnectTimer = null;
    }
    _ws?.close();
    _ws = null;
  }

  function addNotification(payload: { title: string; body: string; type: string }) {
    const n: AppNotification = {
      id: Math.random().toString(36).slice(2),
      title: payload.title,
      body: payload.body,
      type: payload.type,
      read: false,
      timestamp: new Date(),
      route: NOTIFICATION_ROUTES[payload.type] ?? "/",
    };
    notifications.value = [n, ...notifications.value].slice(0, 50);
    toast(payload.title, { description: payload.body });
  }

  function markAsRead(id: string) {
    notifications.value = notifications.value.map((n) =>
      n.id === id ? { ...n, read: true } : n,
    );
  }

  function markAllRead() {
    notifications.value = notifications.value.map((n) => ({ ...n, read: true }));
  }

  function clear() {
    notifications.value = [];
  }

  return {
    notifications,
    unreadCount,
    connect,
    disconnect,
    markAsRead,
    markAllRead,
    clear,
  };
}
