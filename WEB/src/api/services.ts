import { ApiResponse } from "../types/api";
import { Schoolarhisp, SchoolarshipRequestCreate, SchoolarshipRequestUpdate } from "../types/schoolarship";
import api from "./index";

export const getSchoolarship = () => api.get<ApiResponse<Schoolarhisp[]>>("view.schoolarship")
export const addschoolarship = (data: SchoolarshipRequestCreate) => api.post<void>("add.schoolarship", data)
export const updateschoolarship = ( id: number,data: SchoolarshipRequestUpdate) => api.put(`update.schoolarship/${id}`, data)