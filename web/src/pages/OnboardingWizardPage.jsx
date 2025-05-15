import React from "react";
import { useParams } from "react-router-dom";
import OnboardingWizard from "../components/onboarding/OnboardingWizard";
import { useAuth } from "../context/AuthContext";
import Footer from "../components/Footer";

function OnboardingWizardPage() {
  const { step } = useParams();
  const { serverInfo, isLoading } = useAuth();

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col items-center min-h-screen py-8 px-4">
      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl w-full max-w-4xl">
        <h1 className="text-3xl md:text-4xl font-bold mb-8 text-center">
          Getting Started
        </h1>
        {isLoading ? (
          <div className="flex justify-center items-center py-10">
            <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-blue-500"></div>
          </div>
        ) : (
          <OnboardingWizard serverInfo={serverInfo} initialStep={step} />
        )}
      </div>
      
      <Footer />
    </div>
  );
}

export default OnboardingWizardPage;
