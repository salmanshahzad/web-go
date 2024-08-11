import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
    Outlet,
    RouterProvider,
    createRootRoute,
    createRoute,
    createRouter,
    redirect,
} from "@tanstack/react-router";
import axios from "axios";
import { StrictMode } from "react";
import reactDom from "react-dom/client";

import "./index.css";
import { Dashboard } from "./Dashboard";
import { Index } from "./Index";
import { SignIn } from "./SignIn";
import { SignUp } from "./SignUp";

async function isAuthenticated(): Promise<boolean> {
    try {
        await axios.get("/api/user");
        return true;
    } catch {
        return false;
    }
}

const root = document.getElementById("root");
if (!root) {
    throw new Error("Could not find root element");
}

const queryClient = new QueryClient();

const rootRoute = createRootRoute({
    component: () => (
        <>
            <Outlet />
        </>
    ),
});

const dashboardRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/dashboard",
    component: () => <Dashboard />,
    beforeLoad: async () => {
        if (!(await isAuthenticated())) {
            throw redirect({ replace: true, to: "/signin" });
        }
    },
});

const indexRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/",
    component: () => <Index />,
    beforeLoad: async () => {
        if (await isAuthenticated()) {
            throw redirect({ replace: true, to: "/dashboard" });
        }
    },
});

const signInRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/signin",
    component: () => <SignIn />,
    beforeLoad: async () => {
        if (await isAuthenticated()) {
            throw redirect({ replace: true, to: "/dashboard" });
        }
    },
});

const signUpRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/signup",
    component: () => <SignUp />,
    beforeLoad: async () => {
        if (await isAuthenticated()) {
            throw redirect({ replace: true, to: "/dashboard" });
        }
    },
});

const routeTree = rootRoute.addChildren([
    dashboardRoute,
    indexRoute,
    signInRoute,
    signUpRoute,
]);
const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
    interface Register {
        router: typeof router;
    }
}

reactDom.createRoot(root).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <RouterProvider router={router} />
        </QueryClientProvider>
    </StrictMode>,
);
