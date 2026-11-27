import { ApiResponse } from "../types/api";
import { Schoolarhisp } from "../types/schoolarship";
import api from "./index";

export const getSchoolarship = () => api.get<ApiResponse<Schoolarhisp[]>>("view.schoolarship")