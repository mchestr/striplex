import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import PlexSignInButton from "../components/PlexSignInButton";

const LoginPage = () => {
  const [errorMessage, setErrorMessage] = useState("");
  const navigate = useNavigate();
  const { serverInfo, user, isLoading: isAuthLoading, refreshUser } = useAuth();

  // Redirect if user is already authenticated
  useEffect(() => {
    if (user) {
      navigate("/");
    }
  }, [user, navigate]);

  const handleAuthSuccess = () => {
    // Refresh user data and navigate
    refreshUser().then(() => {
      navigate("/");
    });
  };

  const handleAuthError = (message) => {
    setErrorMessage(message);
  };

  if (isAuthLoading) {
    return (
      <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
        <div className="text-xl">Loading...</div>
      </div>
    );
  }

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen overflow-x-hidden">
      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-md w-[90%] relative">
        {/* Ribbon - Updated to include coffee emoji */}
        <a
          href="/stripe/donation"
          className="absolute -right-2 -top-2 overflow-hidden w-28 h-28 z-10 cursor-pointer"
        >
          <div className="absolute transform rotate-45 bg-[#e5a00d]/80 text-[#191a1c] font-bold text-xs py-1 right-[-35px] top-[32px] w-[170px] text-center shadow-sm hover:bg-[#f5b82e]/80 transition-colors duration-200">
            buy me a coffee ☕
          </div>
        </a>

        <h1 className="text-4xl md:text-[3.5rem] font-extrabold mb-4 tracking-tight text-[#f1f2f6] leading-tight">
          {serverInfo.serverName}
        </h1>

        {errorMessage && (
          <div className="mb-6 p-4 bg-red-500/20 border border-red-500/50 rounded-lg text-red-200">
            {errorMessage}
          </div>
        )}

        <div className="space-y-6">
          <PlexSignInButton 
            onSuccess={handleAuthSuccess} 
            onError={handleAuthError} 
            nextUrl="/login-success"
          />
        </div>
      </div>
      <footer className="mt-8 text-sm text-gray-400">
        <p>
          Powered by{" "}
          <a
            href="https://github.com/mchestr/plefi"
            className="text-[#e5a00d] hover:text-[#f5b82e] underline"
            target="_blank"
            rel="noopener noreferrer"
          >
            PleFi
          </a>
        </p>
      </footer>
    </div>
  );
};

export default LoginPage;
