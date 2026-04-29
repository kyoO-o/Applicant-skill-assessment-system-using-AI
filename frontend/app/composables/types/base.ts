export interface BaseModel {
  id: number;
  created_at: string;
  updated_at: string;
}

export interface GetAllResult<T> {
  items: T[];
  total: number;
}
