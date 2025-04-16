import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import BuyMeCoffee from '../components/BuyMeCoffee';

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
        const response = await fetch('/api/v1/stripe/subscriptions');
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
        console.error('Error fetching subscriptions:', error);
        setError("An error occurred while fetching subscriptions");
      } finally {
        setIsLoading(false);
      }
    };

    fetchSubscriptions();
  }, []);

  const formatPrice = (amount, currency = 'USD') => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(amount / 100);
  };

  const formatDate = (timestamp) => {
    return new Date(timestamp * 1000).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const handleCancelSubscription = () => {
    if (!subscriptionToCancel) return;

    fetch('/api/v1/stripe/cancel-subscription', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ subscription_id: subscriptionToCancel.id })
    })
    .then(response => response.json())
    .then(data => {
      if (data.status === "success") {
        const updatedSubscriptions = subscriptions.map(sub => 
          sub.id === subscriptionToCancel.id ? { ...sub, cancel_at_period_end: true } : sub
        );
        setSubscriptions(updatedSubscriptions);
        setShowCancelModal(false);
        setSubscriptionToCancel(null);
      } else {
        alert("Failed to cancel subscription: " + (data.error || "Unknown error"));
      }
    })
    .catch(err => {
      console.error("Error canceling subscription:", err);
      alert("An error occurred while trying to cancel the subscription.");
    });
  };

  const CancelModal = () => (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-[#2d3436] rounded-lg p-6 max-w-md mx-auto">
        <h3 className="text-xl font-bold mb-4">Cancel Subscription</h3>
        <p className="mb-6">
          Are you sure you want to cancel this subscription? You will still have access until the end of your current billing period.
        </p>
        <div className="flex justify-end space-x-4">
          <button
            onClick={() => {
              setShowCancelModal(false);
              setSubscriptionToCancel(null);
            }}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-[#f1f2f6] rounded-md"
          >
            No, Keep Subscription
          </button>
          <button
            onClick={handleCancelSubscription}
            className="px-4 py-2 bg-red-800 hover:bg-red-700 text-red-100 rounded-md"
          >
            Yes, Cancel
          </button>
        </div>
      </div>
    </div>
  );

  if (isLoading) {
    return (
      <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
        <div className="text-xl">Loading subscriptions...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col justify-center items-center min-h-screen">
        <div className="text-xl text-red-500">{error}</div>
        <button 
          onClick={() => navigate('/')} 
          className="mt-4 px-4 py-2 bg-[#e5a00d] hover:bg-[#f5b82e] text-[#191a1c] rounded-md"
        >
          Return Home
        </button>
      </div>
    );
  }

  return (
    <div className="font-sans bg-[#1e272e] text-[#f1f2f6] flex flex-col items-center min-h-screen py-8 px-4">
      {showCancelModal && <CancelModal />}
      <h1 className="text-3xl font-bold mb-8">Your Subscriptions</h1>
      
      <div className="w-full max-w-4xl">
        <div className="overflow-x-auto">
          {subscriptions.length === 0 ? (
            <div className="bg-[#2d3436] rounded-lg p-8 text-center">
              <p className="text-xl mb-4">You don't have any active subscriptions.</p>
            </div>
          ) : (
            <table className="w-full bg-[#2d3436] shadow-lg rounded-lg">
              <thead>
                <tr className="border-b border-gray-700">
                  <th className="px-6 py-4 text-left font-medium">Subscription</th>
                  <th className="px-6 py-4 text-left font-medium">Price</th>
                  <th className="px-6 py-4 text-left font-medium">Status</th>
                  <th className="px-6 py-4 text-left font-medium">Next Billing</th>
                  <th className="px-6 py-4 text-left font-medium">Actions</th>
                </tr>
              </thead>
              <tbody>
                {subscriptions.map((subscription) => (
                  <tr key={subscription.id} className="border-b border-gray-700">
                    <td className="px-6 py-4 flex items-center gap-4">
                      <div className="bg-[#e5a00d] rounded-md w-12 h-12 flex items-center justify-center text-[#191a1c] font-bold text-xl">
                        P
                      </div>
                      <div>
                        <p className="text-[#f1f2f6] font-medium">
                          {subscription.items[0]?.price?.product?.name || "Plex Subscription"}
                        </p>
                        <span className="text-green-400 text-sm">
                          {subscription.cancel_at_period_end ? "Cancels at period end" : "Active"}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      {subscription.items[0]?.price?.unit_amount
                        ? formatPrice(subscription.items[0].price.unit_amount, subscription.items[0].price.currency || 'USD')
                        : "N/A"}
                      {subscription.items[0]?.price?.recurring?.interval && (
                        <span className="text-gray-400 text-sm block">
                          /{subscription.items[0].price.recurring.interval}
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-2 py-1 rounded-full text-xs ${
                        subscription.cancel_at_period_end ? "bg-gray-800 text-gray-300" :
                        subscription.status === "active" ? "bg-green-900 text-green-300" :
                        subscription.status === "past_due" ? "bg-yellow-900 text-yellow-300" :
                        "bg-red-900 text-red-300"
                      }`}>
                        {subscription.cancel_at_period_end ? "cancelled" : subscription.status}
                      </span>
                      {subscription.cancel_at_period_end && (
                        <span className="block text-xs text-gray-400 mt-1">
                          Active until period end
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4">
                      {subscription.items[0]?.current_period_end ? 
                        formatDate(subscription.items[0].current_period_end) : 
                        "N/A"}
                    </td>
                    <td className="px-6 py-4">
                      {!subscription.cancel_at_period_end && (
                        <button
                          onClick={() => {
                            setSubscriptionToCancel(subscription);
                            setShowCancelModal(true);
                          }}
                          className="px-3 py-1 bg-red-800 hover:bg-red-700 text-red-100 rounded-md text-sm"
                        >
                          Cancel
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        
        <div className="mt-6 flex flex-col items-center">
          <div className="mb-4">
            <BuyMeCoffee />
          </div>
          
          <button 
            onClick={() => navigate('/')} 
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-[#f1f2f6] rounded-md"
          >
            Return Home
          </button>
        </div>
      </div>
    </div>
  );
}

export default SubscriptionsPage;
