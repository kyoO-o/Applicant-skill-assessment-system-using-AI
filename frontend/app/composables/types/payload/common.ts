export interface FetchList<T> {
  items: T[];
  total: number;
}

export interface WSEvent {
  type: string;
  data: any;
}
