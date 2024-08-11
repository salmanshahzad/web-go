import { useMutation, useQuery } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import axios from "axios";
import { type ReactNode, useCallback } from "react";

interface User {
    username: string;
}

export function Dashboard(): ReactNode {
    const navigate = useNavigate();
    const signOut = useMutation({
        mutationFn: async () => {
            await axios.delete("/api/session");
        },
        onSuccess: () => {
            navigate({ to: "/", replace: true });
        },
    });
    const onSignOut = useCallback(() => {
        signOut.mutate();
    }, [signOut]);
    const user = useQuery<User>({
        queryKey: ["user"],
        queryFn: async () => {
            const res = await axios.get("/api/user");
            return res.data;
        },
    });

    return (
        <div className="flex gap-4 p-4">
            {user.isLoading && <p>Loading...</p>}
            {user.error && <p>There was an unexpected error.</p>}
            {user.data && (
                <div className="flex gap-4 items-center">
                    <p>Welcome {user.data.username}</p>
                    <button
                        className="btn btn-primary"
                        disabled={signOut.isPending}
                        onClick={onSignOut}
                        type="button"
                    >
                        Sign Out
                    </button>
                </div>
            )}
        </div>
    );
}
