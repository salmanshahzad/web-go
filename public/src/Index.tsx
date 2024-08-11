import { Link } from "@tanstack/react-router";
import type { ReactNode } from "react";

export function Index(): ReactNode {
    return (
        <div className="flex gap-4 p-4">
            <Link className="btn btn-primary" to="/signin">
                Sign In
            </Link>
            <Link className="btn btn-primary" to="/signup">
                Sign Up
            </Link>
        </div>
    );
}
