import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import PlexSignInButton from "../components/PlexSignInButton";
import Footer from "../components/Footer";

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
      
      <Footer />
    </div>
  );
};

export default LoginPage;
