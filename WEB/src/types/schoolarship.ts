export type DiscountType = "percentage" | "amount"
export interface Schoolarhisp {
    id: number
    code: string
    name: string
    discount_type: DiscountType
    discount_percentage?: number | null
    discount_amount?: number | null
    description?: string | null
    active: boolean
}

export interface SchoolarshipRequestCreate {
    code: string
    name: string
    discount_type: DiscountType | null
    discount_percentage?: number | null
    discount_amount?: number | null
    description?: string | null
}

export interface SchoolarshipRequestUpdate {
    code?: string | null
    name?: string | null
    discount_type?: DiscountType | null
    discount_percentage?: number | null
    discount_amount?: number | null
    description?: string | null
    active?: boolean | null
}