import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import axios, { type AxiosError } from "axios";
import type { ReactNode } from "react";
import { Controller, useForm } from "react-hook-form";

import { TextInput } from "./ui/TextInput";
import { type FormErrors, setFormErrors } from "./utils";

interface FormData {
    username: string;
    password: string;
}

export function SignIn(): ReactNode {
    const { control, formState, handleSubmit, setError } = useForm<FormData>({
        defaultValues: {
            username: "",
            password: "",
        },
    });
    const navigate = useNavigate();
    const onSubmit = (data: FormData) => {
        signIn.mutate(data);
    };
    const signIn = useMutation<
        void,
        AxiosError<FormErrors<FormData>>,
        FormData
    >({
        mutationFn: async (data) => {
            await axios.post("/api/session", data);
        },
        onError: (err) => {
            setFormErrors(setError, err);
        },
        onSuccess: () => {
            navigate({ to: "/dashboard", replace: true });
        },
    });

    return (
        <div className="flex flex-col gap-4 max-w-xs p-4">
            <h1 className="font-bold text-3xl">Sign In</h1>
            {formState.errors.root?.message && (
                <div className="alert alert-error">
                    {formState.errors.root.message}
                </div>
            )}
            <form
                className="flex flex-col gap-4"
                onSubmit={handleSubmit(onSubmit)}
            >
                <Controller
                    control={control}
                    name="username"
                    render={({ field, fieldState }) => (
                        <TextInput
                            error={fieldState.error}
                            field={field}
                            label="Username"
                            type="text"
                        />
                    )}
                />
                <Controller
                    control={control}
                    name="password"
                    render={({ field, fieldState }) => (
                        <TextInput
                            error={fieldState.error}
                            field={field}
                            label="Password"
                            type="password"
                        />
                    )}
                />
                <button
                    className="btn btn-primary"
                    disabled={signIn.isPending}
                    type="submit"
                >
                    Sign In
                </button>
            </form>
        </div>
    );
}
