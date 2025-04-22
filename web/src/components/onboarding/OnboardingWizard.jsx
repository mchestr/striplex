import React, { useMemo } from "react";
import { useNavigate, useParams } from "react-router-dom";
import WelcomeStep from "./steps/WelcomeStep";
import RequestsGuideStep from "./steps/RequestsGuideStep";
import DiscordInviteStep from "./steps/DiscordInviteStep";
import TipsAndTricksStep from "./steps/TipsAndTricksStep";

function OnboardingWizard({ initialStep = 0, serverInfo }) {
  // Get the step directly from URL params for more reliable deep linking
  const { step: urlStep } = useParams();
  const urlStepNumber = parseInt(urlStep) || 0;
  const navigate = useNavigate();

  // Dynamically build steps based on available server info
  const steps = useMemo(() => {
    const baseSteps = [{ name: "Welcome", component: WelcomeStep }];

    // Add Requests step only if requestsUrl is provided
    if (serverInfo.requestsUrl) {
      baseSteps.push({ name: "Request Content", component: RequestsGuideStep });
    }

    // Add Discord step if URL is provided
    if (serverInfo.discordServerUrl) {
      baseSteps.push({ name: "Join Community", component: DiscordInviteStep });
    }

    // Always add tips as the final step
    baseSteps.push({ name: "Tips & Tricks", component: TipsAndTricksStep });

    return baseSteps;
  }, [serverInfo]);

  // Calculate the current step safely
  const currentStep = useMemo(() => {
    // Make sure we have steps determined
    if (!steps || steps.length === 0) return 0;

    // Check if URL step is valid
    if (urlStepNumber >= 0 && urlStepNumber < steps.length) {
      return urlStepNumber;
    }

    // If URL step is invalid, redirect to a valid step
    const validStep = Math.max(0, Math.min(initialStep, steps.length - 1));

    // Only redirect if we need to
    if (urlStepNumber !== validStep) {
      // Use setTimeout to avoid navigation during render
      setTimeout(() => {
        navigate(`/onboarding/step/${validStep}`, { replace: true });
      }, 0);
    }

    return validStep;
  }, [steps, urlStepNumber, initialStep, navigate]);

  // Navigation handlers
  const handleNext = () => {
    if (currentStep < steps.length - 1) {
      navigate(`/onboarding/step/${currentStep + 1}`);
    }
  };

  const handlePrev = () => {
    if (currentStep > 0) {
      navigate(`/onboarding/step/${currentStep - 1}`);
    }
  };

  const handleComplete = () => {
    navigate("/");
  };

  // Get the component for the current step
  const CurrentStepComponent =
    steps[currentStep]?.component || steps[0]?.component || WelcomeStep;

  // Determine the next step name for button labeling
  const getNextStepName = () => {
    if (currentStep >= steps.length - 1) {
      return "Complete";
    }
    return steps[currentStep + 1]?.name || "Next";
  };

  return (
    <div className="max-w-3xl mx-auto w-full">
      {/* Progress indicator */}
      <div className="mb-8">
        <div className="flex justify-between">
          {steps.map((step, index) => (
            <div
              key={index}
              className={`text-xs font-medium ${
                currentStep >= index ? "text-blue-400" : "text-gray-500"
              }`}
            >
              {step.name}
            </div>
          ))}
        </div>
        <div className="w-full bg-gray-700 rounded-full h-2.5 mt-2">
          <div
            className="bg-blue-400 h-2.5 rounded-full transition-all duration-300"
            style={{ width: `${((currentStep + 1) / steps.length) * 100}%` }}
          ></div>
        </div>
      </div>

      {/* Step content */}
      <CurrentStepComponent
        onNext={handleNext}
        onPrev={handlePrev}
        onComplete={handleComplete}
        serverInfo={serverInfo}
        nextStepName={getNextStepName()}
        isLastStep={currentStep >= steps.length - 1}
      />
    </div>
  );
}

export default OnboardingWizard;
