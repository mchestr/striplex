import React, { useEffect, useState } from "react";
import Footer from "../components/Footer";

const LoginSuccessPage = () => {
  const [countdown, setCountdown] = useState(5);

  // decrement every second
  useEffect(() => {
    const intervalId = setInterval(
      () => setCountdown((prev) => prev - 1),
      1000
    );
    return () => clearInterval(intervalId);
  }, []);

  // when countdown hits zero, close
  useEffect(() => {
    if (countdown <= 0) handleClose();
  }, [countdown]);

  const handleClose = () => {
    if (window.opener) {
      window.opener.postMessage({ type: "PLEX_AUTH_SUCCESS" }, "*");
    }
    window.close();
  };

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen overflow-x-hidden">
      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-md w-[90%]">
        <div className="mb-6 text-green-400">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-16 w-16 mx-auto"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M5 13l4 4L19 7"
            />
          </svg>
        </div>
        <h1 className="text-3xl font-bold mb-4">Login Successful!</h1>
        <p className="text-lg mb-8">
          You have successfully authenticated with Plex.
        </p>
        <p className="text-sm text-gray-400 mb-6">
          You can now close this window and return to the application.
        </p>

        <button
          onClick={handleClose}
          className="bg-green-500 hover:bg-green-600 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200"
        >
          Close Window
        </button>

        <p className="text-sm mt-2 text-gray-400 mb-4">
          This window will close in {countdown} second
          {countdown !== 1 ? "s" : ""}.
        </p>
      </div>
      
      <Footer />
    </div>
  );
};

export default LoginSuccessPage;
