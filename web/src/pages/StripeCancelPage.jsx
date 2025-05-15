import React from "react";
import { useNavigate } from "react-router-dom";

function StripeCancelPage({ type }) {
  const navigate = useNavigate();

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] min-h-screen py-8 px-4">
      <div className="max-w-3xl mx-auto text-center p-12 rounded-xl shadow-lg bg-[#2d3436] shadow-black/20 w-[90%] md:w-auto">
        <div className="flex items-center justify-center w-20 h-20 mx-auto mb-6 rounded-full text-4xl font-light bg-[#862e2e] text-[#ffc9c9]">
          ✕
        </div>

        <h1 className="text-4xl font-extrabold mb-4">{type} Cancelled</h1>

        <p className="text-lg mb-2 text-[#f1f2f6]">
          Your {type.toLowerCase()} process was cancelled.
        </p>

        <p className="text-lg mb-6 text-[#f1f2f6]">
          If you encountered an issue or have changed your mind, you can try
          again or contact support.
        </p>

        <div className="flex flex-wrap justify-center gap-4 mt-8 md:flex-row flex-col">
          <button
            onClick={() => navigate("/")}
            className="px-7 py-3 bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold rounded-lg shadow-lg hover:shadow-xl transition-all duration-200 text-lg"
          >
            Return Home
          </button>

          {type === "Subscription" ? (
            <button
              onClick={() => (window.location.href = "/stripe/subscribe")}
              className="px-7 py-3 bg-[#34495e] hover:bg-[#2c3e50] text-white font-medium rounded-lg shadow-md hover:shadow-lg transition-all duration-200"
            >
              Try Again
            </button>
          ) : (
            <button
              onClick={() => (window.location.href = "/stripe/donation")}
              className="px-7 py-3 bg-[#34495e] hover:bg-[#2c3e50] text-white font-medium rounded-lg shadow-md hover:shadow-lg transition-all duration-200"
            >
              Try Donating Again
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

export default StripeCancelPage;
