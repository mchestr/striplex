import React, { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import PlexSignInButton from "../components/PlexSignInButton";
import { useAuth } from "../context/AuthContext";
import Footer from "../components/Footer";

function ClaimCodePage() {
  const { user, refreshUser, serverInfo } = useAuth();
  const [code, setCode] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("");
  const [currentStep, setCurrentStep] = useState(user ? 2 : 1); // Step 0: Authentication, Step 1: Claim Code
  const [showSuccessAnimation, setShowSuccessAnimation] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const handleAuthSuccess = () => {
    // Refresh user data and navigate
    refreshUser().then(() => {
      navigate(location.pathname);
    });
  };

  const handleAuthError = (message) => {
    setMessage(message);
  };

  // Extract code from the route path (last segment of URL)
  useEffect(() => {
    // Get the path from location
    const path = location.pathname;
    // Split the path by '/' and get the last segment
    const pathSegments = path.split("/");
    const codeFromPath = pathSegments[pathSegments.length - 1];

    // If the last segment exists and isn't just 'claim-code', set it as the code
    if (codeFromPath && codeFromPath !== "claim") {
      setCode(codeFromPath);
    }
    setCurrentStep(user ? 2 : 1);
  }, [location.pathname, user]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!code.trim()) {
      setMessage("Please enter a valid code");
      setMessageType("error");
      return;
    }

    setIsSubmitting(true);
    setMessage("");

    try {
      const response = await fetch("/api/v1/codes/claim", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ code: code.trim() }),
      });

      const data = await response.json();

      if (response.ok) {
        // Show the success animation
        setShowSuccessAnimation(true);
        // Redirect to the onboarding wizard after animation completes
        setTimeout(() => navigate("/onboarding"), 3000);
      } else {
        setMessage(data.error || "Failed to claim code");
        setMessageType("error");
      }
    } catch (error) {
      console.error("Error claiming code:", error);
      setMessage("An error occurred. Please try again.");
      setMessageType("error");
    } finally {
      setIsSubmitting(false);
    }
  };

  // Success Animation Component
  const SuccessAnimation = () => {
    return (
      <div className="fixed inset-0 bg-black/90 z-50 flex flex-col items-center justify-center overflow-hidden">
        <div className="success-animation">
          <div className="firework-container">
            {[...Array(10)].map((_, i) => (
              <div key={i} className={`firework firework-${i}`}></div>
            ))}
          </div>
          <div className="success-content animate-fade-in">
            <svg
              className="checkmark animate-checkmark"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 52 52"
            >
              <circle
                className="checkmark-circle"
                cx="26"
                cy="26"
                r="25"
                fill="none"
              />
              <path
                className="checkmark-check"
                fill="none"
                d="M14.1 27.2l7.1 7.2 16.7-16.8"
              />
            </svg>
            <h2 className="text-3xl font-bold text-white mt-8 animate-slide-up">
              Code Claimed Successfully!
            </h2>
            <p className="text-xl text-gray-200 mt-4 animate-slide-up animation-delay-300">
              Welcome to {serverInfo?.serverName || "our service"}
            </p>
            <p className="text-lg text-blue-300 mt-6 animate-slide-up animation-delay-500">
              Lets get started with some onboarding!
            </p>
          </div>
        </div>
      </div>
    );
  };

  // Server name display component
  const ServerNameHeader = () => {
    if (!serverInfo || !serverInfo.serverName) return null;

    return (
      <div className="mb-6 text-center">
        <h1 className="text-4xl md:text-[3.5rem] font-extrabold tracking-tight text-[#f1f2f6] leading-tight">
          {serverInfo.serverName}
        </h1>
      </div>
    );
  };

  // Step 1: Authentication required step
  const renderAuthenticationStep = () => {
    return (
      <>
        <ServerNameHeader />
        <h1 className="text-3xl font-bold mb-6">Claim Your Code</h1>

        <div className="bg-[#2a333b] border-l-4 border-[#4b6bfb] p-5 rounded-md mb-6 text-left">
          <h2 className="text-lg font-semibold mb-2 text-gray-100">
            Authentication Required
          </h2>
          <p className="text-gray-300 text-sm">
            Before you can claim your code, you need to sign in with your Plex
            account or create a new one.
          </p>
        </div>

        <div className="flex flex-col space-y-4">
          <PlexSignInButton
            nextUrl={location.pathname}
            onError={handleAuthError}
            onSuccess={handleAuthSuccess}
          />
          <p className="text-amber-300/80 text-xs mt-1 flex items-center justify-center">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-3 w-3 mr-1"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            Unverified app - Plex may show a warning during sign-in
          </p>
        </div>
      </>
    );
  };

  // Step 2: Claim code step
  const renderClaimCodeStep = () => {
    return (
      <>
        <ServerNameHeader />

        {message && (
          <div
            className={`p-4 mb-6 rounded-lg ${
              messageType === "success"
                ? "bg-green-500/20 border border-green-500/50 text-green-200"
                : "bg-red-500/20 border border-red-500/50 text-red-200"
            }`}
          >
            {message}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-6 relative">
            <label
              htmlFor="code"
              className="block text-sm font-medium text-gray-300 mb-1"
            >
              Enter Invite Code
            </label>
            <div className="relative">
              <input
                id="code"
                type="text"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                className="w-full p-3 pr-12 bg-[#3a4149] border border-gray-600 rounded-lg text-white focus:ring-blue-500 focus:border-blue-500"
                disabled={isSubmitting}
              />
              <button
                type="submit"
                disabled={isSubmitting}
                className={`absolute right-2 top-1/2 transform -translate-y-1/2 p-2 rounded-full ${
                  isSubmitting
                    ? "bg-[#3a4149] text-gray-500"
                    : "bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c]"
                } transition-all duration-200 disabled:opacity-70 disabled:cursor-not-allowed`}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-5 w-5"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fillRule="evenodd"
                    d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z"
                    clipRule="evenodd"
                  />
                </svg>
              </button>
            </div>
          </div>
        </form>
      </>
    );
  };

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
      {showSuccessAnimation && <SuccessAnimation />}
      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-md w-[90%]">
        {currentStep === 1 ? renderAuthenticationStep() : renderClaimCodeStep()}

        <div className="mt-8">
          <button
            onClick={() => navigate("/")}
            className="text-sm text-gray-400 hover:text-gray-200 transition-colors"
          >
            Back to Home
          </button>
        </div>
      </div>
      
      <Footer />
      <style jsx>{`
        @keyframes scale-in {
          0% {
            transform: scale(0);
            opacity: 0;
          }
          100% {
            transform: scale(1);
            opacity: 1;
          }
        }

        @keyframes draw-check {
          0% {
            stroke-dashoffset: 100;
          }
          100% {
            stroke-dashoffset: 0;
          }
        }

        @keyframes draw-circle {
          0% {
            stroke-dashoffset: 600;
          }
          100% {
            stroke-dashoffset: 0;
          }
        }

        @keyframes firework-animation {
          0% {
            transform: translate(0, 0);
            width: 0px;
            height: 0px;
            opacity: 1;
          }
          100% {
            transform: translate(var(--x), var(--y));
            width: var(--size);
            height: var(--size);
            opacity: 0;
          }
        }

        .checkmark {
          width: 100px;
          height: 100px;
          border-radius: 50%;
          display: block;
          margin: 0 auto;
        }

        .checkmark-circle {
          stroke-width: 2;
          stroke: #4bb71b;
          stroke-dasharray: 166;
          stroke-dashoffset: 166;
          fill: none;
          animation: draw-circle 1s cubic-bezier(0.65, 0, 0.45, 1) forwards 0.5s;
        }

        .checkmark-check {
          stroke-width: 2;
          stroke: #4bb71b;
          stroke-dasharray: 48;
          stroke-dashoffset: 48;
          fill: none;
          animation: draw-check 0.8s cubic-bezier(0.65, 0, 0.45, 1) forwards 1s;
        }

        .success-animation {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          height: 100%;
        }

        .animate-fade-in {
          animation: scale-in 0.5s ease forwards;
        }

        .animate-slide-up {
          transform: translateY(20px);
          opacity: 0;
          animation: slide-up 0.6s ease forwards;
        }

        .animation-delay-300 {
          animation-delay: 0.3s;
        }

        .animation-delay-500 {
          animation-delay: 0.5s;
        }

        @keyframes slide-up {
          0% {
            transform: translateY(20px);
            opacity: 0;
          }
          100% {
            transform: translateY(0);
            opacity: 1;
          }
        }

        .firework-container {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          pointer-events: none;
        }

        .firework {
          position: absolute;
          top: 50%;
          left: 50%;
          border-radius: 50%;
          background-image: radial-gradient(
            circle,
            #ff0080,
            #ff8c00,
            #40e0d0,
            #7fff00,
            #ff0080
          );
          background-size: 600% 600%;
          background-position: center center;
          animation: firework-animation 1.5s ease-out forwards;
          mix-blend-mode: screen;
        }

        .firework-0 {
          --x: -180px;
          --y: -200px;
          --size: 100px;
          animation-delay: 0.2s;
        }
        .firework-1 {
          --x: 200px;
          --y: -150px;
          --size: 125px;
          animation-delay: 0.4s;
        }
        .firework-2 {
          --x: -130px;
          --y: 180px;
          --size: 80px;
          animation-delay: 0.6s;
        }
        .firework-3 {
          --x: 150px;
          --y: 180px;
          --size: 110px;
          animation-delay: 0.8s;
        }
        .firework-4 {
          --x: -80px;
          --y: -120px;
          --size: 90px;
          animation-delay: 0.3s;
        }
        .firework-5 {
          --x: 100px;
          --y: -80px;
          --size: 130px;
          animation-delay: 0.5s;
        }
        .firework-6 {
          --x: -200px;
          --y: 70px;
          --size: 70px;
          animation-delay: 0.7s;
        }
        .firework-7 {
          --x: 180px;
          --y: 100px;
          --size: 85px;
          animation-delay: 0.9s;
        }
        .firework-8 {
          --x: -50px;
          --y: -180px;
          --size: 120px;
          animation-delay: 0.2s;
        }
        .firework-9 {
          --x: 30px;
          --y: 150px;
          --size: 95px;
          animation-delay: 0.4s;
        }
      `}</style>
    </div>
  );
}

export default ClaimCodePage;
