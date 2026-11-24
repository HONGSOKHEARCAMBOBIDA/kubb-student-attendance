import { useRouter } from "vue-router";

export const useBack = () => {
    const router = useRouter();

    const goback = (): void => {
        router.back();
    };

    return {
        goback,
    };
};