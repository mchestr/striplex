import React, {
  createContext,
  useState,
  useEffect,
  useContext,
  useMemo,
} from "react";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [serverInfo, setServerInfo] = useState({
    serverName: "PleFi",
    discordServerUrl: "",
    requestsUrl: "",
    features: [],
  });
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchInitialData = async () => {
      try {
        const serverInfoResponse = await fetch("/info");
        if (serverInfoResponse.ok) {
          const serverData = await serverInfoResponse.json();
          setServerInfo({
            serverName: serverData.server_name ?? "PleFi",
            discordServerUrl: serverData.discord_server_url ?? "",
            requestsUrl: serverData.requests_url ?? "",
            features: serverData.features ?? [],
          });
        }

        const userResponse = await fetch("/api/v1/user/me");
        if (userResponse.ok) {
          const userData = await userResponse.json();
          if (userData.status === "success" && userData.user) {
            setUser(userData.user);
          }
        }
      } catch (error) {
        console.error("Error fetching initial app data:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchInitialData();
  }, []);

  const refreshUser = async () => {
    try {
      const userResponse = await fetch("/api/v1/user/me");
      if (userResponse.ok) {
        const userData = await userResponse.json();
        if (userData.status === "success" && userData.user) {
          setUser(userData.user);
          return userData.user;
        }
      }
      return null;
    } catch (error) {
      console.error("Error refreshing user data:", error);
      return null;
    }
  };

  const logout = async () => {
    try {
      await fetch("/logout", { method: "POST" });
      setUser(null);
    } catch (error) {
      console.error("Error during logout:", error);
    }
  };

  const value = useMemo(
    () => ({
      serverInfo,
      user,
      isLoading,
      isAuthenticated: !!user,
      refreshUser,
      logout,
    }),
    [serverInfo, user, isLoading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export default AuthContext;
