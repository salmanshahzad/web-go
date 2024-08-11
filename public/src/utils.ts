import type { AxiosError } from "axios";
import type { FieldValues, Path, UseFormSetError } from "react-hook-form";

export interface FormErrors<T extends FieldValues> {
    errors?: Record<Path<T>, string[]>;
    message?: string;
}

export function setFormErrors<T extends FieldValues>(
    setError: UseFormSetError<T>,
    err: AxiosError<FormErrors<T>>,
) {
    const data = err.response?.data;
    if (!data) {
        setError("root", { message: "There was an unexpected error" });
    } else if (data.errors) {
        const keys = Object.keys(data.errors) as Path<T>[];
        for (const key of keys) {
            const message = data.errors[key][0];
            if (message) {
                setError(key, { message });
            }
        }
    } else if (data.message) {
        setError("root", { message: data.message });
    }
}
