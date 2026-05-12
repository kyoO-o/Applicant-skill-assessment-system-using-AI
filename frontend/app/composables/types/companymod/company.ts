import type { BaseModel } from "../base";

export interface Company extends BaseModel {
  name: string;
  description: string;
  register_id: string;
  city: string;
  district: string;
  location_x: number;
  location_y: number;
}
