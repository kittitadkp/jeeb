import { createContext } from "react";

export interface AuthUser {
  id: string;
  email: string;
  name: string;
}

export interface AuthContextValue {
  user: AuthUser | null;
  ready: boolean;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextValue>({
  user: null,
  ready: false,
  logout: () => {},
});
