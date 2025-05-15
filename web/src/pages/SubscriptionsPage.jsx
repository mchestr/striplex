import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Footer from "../components/Footer";

function SubscriptionsPage() {
  const [isLoading, setIsLoading] = useState(true);
  const [subscriptions, setSubscriptions] = useState([]);
  const [error, setError] = useState(null);
  const [showCancelModal, setShowCancelModal] = useState(false);
  const [subscriptionToCancel, setSubscriptionToCancel] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchSubscriptions = async () => {
      try {
        const response = await fetch("/api/v1/stripe/subscriptions");
        if (response.ok) {
          const data = await response.json();
          if (data.status === "success") {
            setSubscriptions(data.subscriptions || []);
          } else {
            setError("Failed to load subscriptions data");
          }
        } else {
          setError("Failed to load subscriptions");
        }
      } catch (error) {
        console.error("Error fetching subscriptions:", error);
        setError("An error occurred while fetching subscriptions");
      } finally {
        setIsLoading(false);
      }
    };

    fetchSubscriptions();
  }, []);

  const formatPrice = (amount, currency = "USD") => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: currency,
    }).format(amount / 100);
  };

  const formatDate = (timestamp) => {
    return new Date(timestamp * 1000).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const handleCancelSubscription = () => {
    if (!subscriptionToCancel) return;

    fetch("/api/v1/stripe/cancel-subscription", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ subscription_id: subscriptionToCancel.id }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.status === "success") {
          const updatedSubscriptions = subscriptions.map((sub) =>
            sub.id === subscriptionToCancel.id
              ? { ...sub, cancel_at_period_end: true }
              : sub
          );
          setSubscriptions(updatedSubscriptions);
          setShowCancelModal(false);
          setSubscriptionToCancel(null);
        } else {
          alert(
            "Failed to cancel subscription: " + (data.error || "Unknown error")
          );
        }
      })
      .catch((err) => {
        console.error("Error canceling subscription:", err);
        alert("An error occurred while trying to cancel the subscription.");
      });
  };

  if (isLoading) {
    return (
      <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
        <div className="flex flex-col items-center justify-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[#e5a00d] mb-4"></div>
          <div className="text-xl">Loading subscriptions...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
        <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-md w-[90%]">
          <div className="text-3xl text-red-500 mb-2">Oops!</div>
          <div className="text-xl text-red-400 mb-6">{error}</div>
          <button
            onClick={() => navigate("/")}
            className="px-6 py-3 bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold rounded-lg shadow-lg hover:shadow-xl transition-all duration-200 text-lg"
          >
            Return Home
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen overflow-x-hidden">
      {showCancelModal && (
        <CancelModal
          setShowCancelModal={setShowCancelModal}
          setSubscriptionToCancel={setSubscriptionToCancel}
          handleCancelSubscription={handleCancelSubscription}
        />
      )}

      <div className="bg-[#2d3436] shadow-lg shadow-black/20 p-8 md:p-12 rounded-xl text-center max-w-4xl w-[90%]">
        <h1 className="text-4xl md:text-[3.5rem] font-extrabold mb-2 tracking-tight text-[#f1f2f6] leading-tight">
          Your Subscriptions
        </h1>
        <p className="text-gray-400 text-center mb-8">
          Manage your active subscriptions and billing
        </p>

        {subscriptions.length === 0 ? (
          <div className="bg-[#34495e] rounded-lg p-8 text-center mb-6">
            <div className="flex justify-center mb-4">
              <div className="bg-gray-700 p-4 rounded-full">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-16 w-16 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={1.5}
                    d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"
                  />
                </svg>
              </div>
            </div>
            <h2 className="text-2xl font-bold mb-2">No Active Subscriptions</h2>
            <p className="text-gray-300 mb-6">
              You don't have any active subscriptions at the moment.
            </p>
            <button
              onClick={() => (window.location.href = "/stripe/subscribe")}
              className="px-6 py-3 bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold rounded-lg shadow-lg hover:shadow-xl transition-all duration-200 text-lg"
            >
              Subscribe Now
            </button>
          </div>
        ) : (
          <div className="space-y-6 mb-6">
            {subscriptions.map((subscription) => (
              <div
                key={subscription.id}
                className="bg-[#34495e] shadow-md rounded-lg overflow-hidden"
              >
                {/* Subscription Header - Status Bar */}
                <div
                  className={`h-2 w-full ${
                    subscription.cancel_at_period_end
                      ? "bg-gray-600"
                      : subscription.status === "active" ||
                        subscription.status === "trialing"
                      ? "bg-green-500"
                      : subscription.status === "past_due"
                      ? "bg-yellow-500"
                      : "bg-red-500"
                  }`}
                ></div>

                <div className="p-6">
                  {/* Subscription Main Content */}
                  <div className="flex flex-col md:flex-row items-start md:items-center">
                    {/* Logo/Icon */}
                    <div className="bg-[#e5a00d] rounded-lg w-16 h-16 flex items-center justify-center text-[#191a1c] font-bold text-2xl mb-4 md:mb-0 md:mr-6">
                      P
                    </div>

                    {/* Subscription Info */}
                    <div className="flex-grow">
                      <div className="flex flex-col md:flex-row md:items-center md:justify-between">
                        <div>
                          <h2 className="text-2xl font-bold">
                            {subscription.items[0]?.price?.product?.name ||
                              "Plex Subscription"}
                          </h2>
                          <div className="flex items-center mt-1">
                            <span
                              className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium mr-2 ${
                                subscription.cancel_at_period_end
                                  ? "bg-gray-800 text-gray-300"
                                  : subscription.status === "active" ||
                                    subscription.status === "trialing"
                                  ? "bg-green-900 text-green-300"
                                  : subscription.status === "past_due"
                                  ? "bg-yellow-900 text-yellow-300"
                                  : "bg-red-900 text-red-300"
                              }`}
                            >
                              <span className="h-2 w-2 mr-1 rounded-full bg-current"></span>
                              {subscription.cancel_at_period_end
                                ? "Cancels soon"
                                : subscription.status === "active" ||
                                  subscription.status === "trialing"
                                ? "Active"
                                : subscription.status}
                            </span>

                            {subscription.cancel_at_period_end && (
                              <span className="text-sm text-gray-400">
                                Active until period end
                              </span>
                            )}
                          </div>
                        </div>

                        <div className="mt-4 md:mt-0 text-right">
                          <div className="text-2xl font-bold">
                            {subscription.items[0]?.price?.unit_amount
                              ? formatPrice(
                                  subscription.items[0].price.unit_amount,
                                  subscription.items[0].price.currency || "USD"
                                )
                              : "N/A"}
                          </div>
                          {subscription.items[0]?.price?.recurring
                            ?.interval && (
                            <div className="text-gray-400 text-sm">
                              per{" "}
                              {subscription.items[0].price.recurring.interval}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Subscription Details */}
                  <div className="mt-6 grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="bg-[#2d3436] p-4 rounded-lg">
                      <h3 className="text-gray-400 text-sm font-medium mb-1">
                        Billing Period
                      </h3>
                      <div className="text-lg font-medium">
                        {subscription.items[0]?.price?.recurring
                          ?.interval_count || 1}{" "}
                        {subscription.items[0]?.price?.recurring?.interval ||
                          "month"}
                      </div>
                    </div>

                    <div className="bg-[#2d3436] p-4 rounded-lg">
                      <h3 className="text-gray-400 text-sm font-medium mb-1">
                        Next Billing Date
                      </h3>
                      <div className="text-lg font-medium">
                        {subscription.items[0]?.current_period_end
                          ? formatDate(subscription.items[0].current_period_end)
                          : "N/A"}
                      </div>
                    </div>

                    <div className="bg-[#2d3436] p-4 rounded-lg">
                      <h3 className="text-gray-400 text-sm font-medium mb-1">
                        Payment Method
                      </h3>
                      <div className="text-lg font-medium flex items-center">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-5 w-5 mr-2 text-gray-400"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"
                          />
                        </svg>
                        Stripe
                      </div>
                    </div>
                  </div>

                  {/* Subscription Actions */}
                  <div className="mt-6 pt-4 border-t border-gray-700 flex justify-end space-x-4">
                    {!subscription.cancel_at_period_end ? (
                      <button
                        onClick={() => {
                          setSubscriptionToCancel(subscription);
                          setShowCancelModal(true);
                        }}
                        className="px-4 py-2 bg-red-800 hover:bg-red-700 text-red-100 rounded-lg transition-all duration-200 flex items-center"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-5 w-5 mr-2"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M6 18L18 6M6 6l12 12"
                          />
                        </svg>
                        Cancel Subscription
                      </button>
                    ) : (
                      <button
                        onClick={() =>
                          (window.location.href = "/stripe/subscribe")
                        }
                        className="px-4 py-2 bg-blue-700 hover:bg-blue-600 text-white rounded-lg transition-all duration-200 flex items-center"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          className="h-5 w-5 mr-2"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                          />
                        </svg>
                        Renew Subscription
                      </button>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        <div className="mt-6">
          <button
            onClick={() => navigate("/")}
            className="px-6 py-3 bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] font-bold rounded-lg shadow-lg hover:shadow-xl transition-all duration-200 text-lg"
          >
            Return Home
          </button>
        </div>
      </div>

      <Footer />
    </div>
  );
}

const CancelModal = ({
  setShowCancelModal,
  setSubscriptionToCancel,
  handleCancelSubscription,
}) => (
  <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50 p-4">
    <div className="bg-[#2d3436] rounded-xl p-8 max-w-md w-[90%] shadow-xl">
      <div className="mb-2 flex justify-center">
        <div className="bg-red-900/40 p-4 rounded-full">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-10 w-10 text-red-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1.5}
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            />
          </svg>
        </div>
      </div>
      <h3 className="text-2xl font-bold mb-2 text-center">
        Cancel Subscription
      </h3>
      <p className="mb-6 text-gray-300 text-center">
        Are you sure you want to cancel this subscription? You will still have
        access until the end of your current billing period.
      </p>
      <div className="flex flex-col space-y-3">
        <button
          onClick={handleCancelSubscription}
          className="w-full px-4 py-3 bg-red-800 hover:bg-red-700 text-red-100 rounded-lg transition-all duration-200"
        >
          Yes, Cancel Subscription
        </button>
        <button
          onClick={() => {
            setShowCancelModal(false);
            setSubscriptionToCancel(null);
          }}
          className="w-full px-4 py-3 bg-gray-700 hover:bg-gray-600 text-[#f1f2f6] rounded-lg transition-all duration-200"
        >
          No, Keep Subscription
        </button>
      </div>
    </div>
  </div>
);

export default SubscriptionsPage;
