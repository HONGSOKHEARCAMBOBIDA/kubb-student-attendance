export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface PaginationMetadata {
    page: number
    pageSize: number
    totalCount: number
    totalPages: number
    hasNext: boolean
    hasPrev: boolean
    current: number
}

export interface PaginatedResponse<T> {
    data: T[]
    pagination: PaginationMetadata
}