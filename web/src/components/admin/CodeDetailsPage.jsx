import React, { useState, useEffect } from 'react';

function CodeDetailsPage({ codeId, onBack }) {
  const [codeDetails, setCodeDetails] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchCodeDetails = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(`/api/v1/codes/${codeId}`);
        if (!response.ok) {
          throw new Error(`Failed to fetch code details: ${response.status}`);
        }
        const data = await response.json();
        setCodeDetails(data);
      } catch (error) {
        console.error('Error fetching code details:', error);
        setError(error.message);
      } finally {
        setIsLoading(false);
      }
    };

    if (codeId) {
      fetchCodeDetails();
    }
  }, [codeId]);

  const formatDate = (timestamp) => {
    if (!timestamp) return 'N/A';
    const date = new Date(timestamp);
    return date.toLocaleString();
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-xl">Loading code details...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6 bg-[#2d3436] rounded-lg">
        <div className="flex items-center mb-6">
          <button
            onClick={onBack}
            className="mr-4 p-2 bg-gray-700 hover:bg-gray-600 rounded-lg"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
          </button>
          <h2 className="text-2xl font-bold">Code Details</h2>
        </div>

        <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-4 text-red-200">
          {error}
        </div>
      </div>
    );
  }

  if (!codeDetails) {
    return (
      <div className="p-6 bg-[#2d3436] rounded-lg">
        <div className="flex items-center mb-6">
          <button
            onClick={onBack}
            className="mr-4 p-2 bg-gray-700 hover:bg-gray-600 rounded-lg"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
          </button>
          <h2 className="text-2xl font-bold">Code Details</h2>
        </div>

        <div className="bg-gray-700/50 p-4 rounded-lg">
          No code details found
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="p-6 bg-[#2d3436] rounded-lg">
        <div className="flex items-center mb-6">
          <button
            onClick={onBack}
            className="mr-4 p-2 bg-gray-700 hover:bg-gray-600 rounded-lg"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
          </button>
          <h2 className="text-2xl font-bold">Code Details</h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-[#1e272e] p-5 rounded-lg">
            <div className="mb-4">
              <h3 className="text-sm font-medium text-gray-400 mb-1">Invite Code</h3>
              <p className="text-xl font-mono bg-gray-800 p-3 rounded">{codeDetails.code}</p>
            </div>
            <div className="mb-4">
              <h3 className="text-sm font-medium text-gray-400 mb-1">Status</h3>
              <p className="flex items-center">
                <span className={`inline-block w-3 h-3 rounded-full mr-2 ${
                  codeDetails.is_disabled
                    ? "bg-gray-500"
                    : codeDetails.used_count >= codeDetails.max_uses
                      ? "bg-gray-500"
                      : new Date(codeDetails.expires_at) < new Date()
                        ? "bg-red-500"
                        : "bg-green-500"
                }`}></span>
                <span>
                  {codeDetails.is_disabled
                    ? "Disabled"
                    : codeDetails.used_count >= codeDetails.max_uses
                      ? "Used"
                      : new Date(codeDetails.expires_at) < new Date()
                        ? "Expired"
                        : "Active"
                  }
                </span>
              </p>
            </div>
            <div className="mb-4">
              <h3 className="text-sm font-medium text-gray-400 mb-1">Usage</h3>
              <p>{codeDetails.used_count} / {codeDetails.max_uses}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-400 mb-1">Entitlement</h3>
              <p>{codeDetails.entitlement_name || "N/A"}</p>
            </div>
          </div>

          <div className="bg-[#1e272e] p-5 rounded-lg">
            <div className="mb-4">
              <h3 className="text-sm font-medium text-gray-400 mb-1">Created At</h3>
              <p>{formatDate(codeDetails.created_at)}</p>
            </div>
            <div className="mb-4">
              <h3 className="text-sm font-medium text-gray-400 mb-1">Expires At</h3>
              <p>{formatDate(codeDetails.expires_at)}</p>
            </div>
            <div className="mb-4">
              <h3 className="text-sm font-medium text-gray-400 mb-1">Duration</h3>
              <p>{formatDate(codeDetails.duration)}</p>
            </div>
            {codeDetails.created_by && (
              <div>
                <h3 className="text-sm font-medium text-gray-400 mb-1">Created By</h3>
                <p>{codeDetails.created_by}</p>
              </div>
            )}
          </div>
        </div>
      </div>

      {codeDetails.users && codeDetails.users.length > 0 && (
        <div className="p-6 bg-[#2d3436] rounded-lg">
          <h2 className="text-xl font-bold mb-4">Users Who Redeemed This Code</h2>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-[#1e272e]/50 text-left">
                <tr className="border-b border-gray-700">
                  <th scope="col" className="px-6 py-3 font-medium">Username</th>
                  <th scope="col" className="px-6 py-3 font-medium">Email</th>
                  <th scope="col" className="px-6 py-3 font-medium">Status</th>
                  <th scope="col" className="px-6 py-3 font-medium">Date</th>
                </tr>
              </thead>
              <tbody>
                {codeDetails.users.map((user) => (
                  <tr key={user.id} className="border-b border-gray-700">
                    <td className="px-6 py-4">
                      <a 
                        href={`/admin/users/${user.id}`}
                        className="text-blue-400 hover:text-blue-300 hover:underline"
                      >
                        {user.username}
                      </a>
                    </td>
                    <td className="px-6 py-4">{user.email}</td>
                    <td className="px-6 py-4">
                      {user.is_admin ? (
                        <span className="bg-green-900/30 text-green-400 px-2 py-1 rounded-full text-xs">
                          Admin
                        </span>
                      ) : (
                        <span className="bg-blue-900/30 text-blue-400 px-2 py-1 rounded-full text-xs">
                          User
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4">{formatDate(user.created_at)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}

export default CodeDetailsPage;
