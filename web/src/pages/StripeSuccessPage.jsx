import React from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

function StripeSuccessPage({ type }) {
  const navigate = useNavigate();
  const { user } = useAuth();

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] min-h-screen py-8 px-4">
      <div className="max-w-3xl mx-auto text-center p-12 rounded-xl shadow-lg bg-[#2d3436] shadow-black/20">
        <div className={`flex items-center justify-center w-20 h-20 mx-auto mb-6 rounded-full text-4xl font-light bg-[#2b8a3e] text-[#e3f9e5]`}>
          {type === "Donation" ? "❤️" : "✓"}
        </div>

        <h1 className="text-4xl font-extrabold mb-4">
          {type === "Donation" ? "Thank You for Your Donation!" : `${type} Successful!`}
        </h1>

        {type === "Subscription" ? (
          <p className="text-lg mb-6 text-[#f1f2f6]">
            Thank you{user ? `, ${user.username}` : ""}! Your
            subscription will be activated soon.
          </p>
        ) : (
          <p className="text-lg mb-6 text-[#f1f2f6]">
            Your support is greatly appreciated{user ? `, ${user.username}` : ""}!
          </p>
        )}

        {type === "Donation" ? (
          <p className="text-lg mb-6 text-[#f1f2f6]">
            We <span className="inline-block animate-pulse">❤️</span> supporters
            like you!
          </p>
        ) : (
          <>
            <p className="text-lg mb-6 text-[#f1f2f6]">
              An invite has been sent to your account and should have been auto-accepted.
            </p>
            <p className="text-lg mb-6 text-[#f1f2f6]">
              You may need to accept the invite within your Plex account to gain
              full access.
            </p>
          </>
        )}

        <div className="flex flex-wrap justify-center gap-4 mt-8">
          {type === "Donation" ? (
            <button
              onClick={() => navigate("/")}
              className="px-7 py-3 bg-[#4b6bfb] hover:bg-[#3557fa] text-white font-bold rounded-lg shadow-md hover:shadow-lg transition transform hover:-translate-y-0.5"
            >
              Return Home
            </button>
          ) : (
            <>
              <button
                onClick={() => navigate("/onboarding")}
                className="px-7 py-3 bg-[#4b6bfb] hover:bg-[#3557fa] text-white font-bold rounded-lg shadow-md hover:shadow-lg transition transform hover:-translate-y-0.5"
              >
                Get Started
              </button>

              <a
                href="https://app.plex.tv/desktop/#!/settings/manage-library-access"
                target="_blank"
                rel="noopener noreferrer"
                className="px-7 py-3 bg-[#ffb142] hover:bg-[#ff9f1a] text-white font-bold rounded-lg shadow-md hover:shadow-lg transition transform hover:-translate-y-0.5"
              >
                Check Plex Requests
              </a>
            </>
          )}
        </div>
      </div>
    </div>
  );
}

export default StripeSuccessPage;
