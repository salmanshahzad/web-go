import classNames from "classnames";
import type { InputHTMLAttributes, ReactNode } from "react";
import type {
    ControllerRenderProps,
    FieldError,
    FieldPath,
    FieldValues,
} from "react-hook-form";

export interface TextInputProps<T extends FieldValues, U extends FieldPath<T>>
    extends InputHTMLAttributes<HTMLInputElement> {
    error: FieldError | undefined;
    field: ControllerRenderProps<T, U>;
    label: string;
}

export function TextInput<T extends FieldValues, U extends FieldPath<T>>(
    props: TextInputProps<T, U>,
): ReactNode {
    const { className, error, field, label, ...rest } = props;

    return (
        <label className="form-control w-full">
            <div className="label">
                <span className="label-text">{label}</span>
            </div>
            <input
                className={classNames(
                    "input",
                    "input-bordered",
                    "w-full",
                    {
                        "input-error": error,
                    },
                    className,
                )}
                {...field}
                {...rest}
            />
            <div className="label">
                <span className="label-text-alt text-red-400">
                    {error?.message}
                </span>
            </div>
        </label>
    );
}
