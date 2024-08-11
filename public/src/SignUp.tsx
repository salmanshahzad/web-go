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
    confirmPassword: string;
}

export function SignUp(): ReactNode {
    const { control, formState, handleSubmit, setError } = useForm<FormData>({
        defaultValues: {
            username: "",
            password: "",
            confirmPassword: "",
        },
    });
    const navigate = useNavigate();
    const onSubmit = (data: FormData) => {
        if (data.password !== data.confirmPassword) {
            setError("root", {
                message: "Passwords do not match",
            });
            return;
        }
        signUp.mutate(data);
    };
    const signUp = useMutation<
        void,
        AxiosError<FormErrors<FormData>>,
        FormData
    >({
        mutationFn: async (data) => {
            await axios.post("/api/user", data);
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
            <h1 className="font-bold text-3xl">Sign Up</h1>
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
                <Controller
                    control={control}
                    name="confirmPassword"
                    render={({ field, fieldState }) => (
                        <TextInput
                            error={fieldState.error}
                            field={field}
                            label="Confirm Password"
                            type="password"
                        />
                    )}
                />
                <button
                    className="btn btn-primary"
                    disabled={signUp.isPending}
                    type="submit"
                >
                    Sign Up
                </button>
            </form>
        </div>
    );
}
