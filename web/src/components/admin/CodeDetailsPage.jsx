import React, { useState, useEffect } from "react";

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
        console.error("Error fetching code details:", error);
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
    if (!timestamp) return "N/A";
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
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
              ></path>
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
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M10 19l-7-7m0 0l7-7m-7 7h18"
              ></path>
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
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center">
            <button
              onClick={onBack}
              className="mr-4 p-2 bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M10 19l-7-7m0 0l7-7m-7 7h18"
                ></path>
              </svg>
            </button>
            <h2 className="text-2xl font-bold flex items-center">
              <svg
                className="w-6 h-6 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                ></path>
              </svg>
              Invite Code Details
            </h2>
          </div>

          <div
            className={`px-3 py-1 rounded-md text-sm font-medium ${
              codeDetails.is_disabled
                ? "bg-gray-700/50 text-gray-300 border border-gray-600"
                : codeDetails.used_count >= codeDetails.max_uses
                ? "bg-gray-600/50 text-gray-300 border border-gray-500"
                : new Date(codeDetails.expires_at) < new Date()
                ? "bg-red-900/30 text-red-300 border border-red-700"
                : "bg-green-900/30 text-green-300 border border-green-700"
            }`}
          >
            {codeDetails.is_disabled
              ? "Disabled"
              : codeDetails.used_count >= codeDetails.max_uses
              ? "Fully Used"
              : new Date(codeDetails.expires_at) < new Date()
              ? "Expired"
              : "Active"}
          </div>
        </div>

        <div className="mb-6">
          <div className="bg-[#1e272e] p-4 rounded-lg flex flex-col md:flex-row md:items-center md:space-x-6">
            <div className="mb-4 md:mb-0 flex-grow">
              <p className="text-xs uppercase tracking-wider text-blue-400 font-semibold mb-1">
                Invite Code
              </p>
              <div className="flex items-center">
                <code className="text-xl font-mono bg-gray-800 p-3 rounded-md flex-grow">
                  {codeDetails.code}
                </code>
              </div>
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div className="text-center p-3 bg-[#2d3436] rounded-md">
                <p className="text-xs uppercase tracking-wider text-gray-400 font-semibold mb-1">
                  Usage
                </p>
                <p className="text-xl font-semibold">
                  {codeDetails.used_count}/{codeDetails.max_uses}
                </p>
              </div>

              <div className="text-center p-3 bg-[#2d3436] rounded-md">
                <p className="text-xs uppercase tracking-wider text-gray-400 font-semibold mb-1">
                  Status
                </p>
                <p className="font-medium">
                  <span
                    className={`inline-block w-3 h-3 rounded-full ${
                      codeDetails.is_disabled
                        ? "bg-gray-500"
                        : codeDetails.used_count >= codeDetails.max_uses
                        ? "bg-gray-500"
                        : new Date(codeDetails.expires_at) < new Date()
                        ? "bg-red-500"
                        : "bg-green-500"
                    } mr-2`}
                  ></span>
                  {codeDetails.is_disabled
                    ? "Disabled"
                    : codeDetails.used_count >= codeDetails.max_uses
                    ? "Used"
                    : new Date(codeDetails.expires_at) < new Date()
                    ? "Expired"
                    : "Active"}
                </p>
              </div>

              <div className="text-center p-3 bg-[#2d3436] rounded-md">
                <p className="text-xs uppercase tracking-wider text-gray-400 font-semibold mb-1">
                  Entitlement
                </p>
                <p className="truncate font-medium">
                  {codeDetails.entitlement_name || "N/A"}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-[#1e272e] p-5 rounded-lg">
            <h3 className="text-sm uppercase tracking-wider text-blue-400 font-semibold mb-3 border-b border-gray-700 pb-2">
              Time Information
            </h3>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-xs text-gray-400 font-medium mb-1">
                  Created
                </p>
                <p className="font-medium">{formatDate(codeDetails.created_at)}</p>
              </div>

              <div>
                <p className="text-xs text-gray-400 font-medium mb-1">
                  Expires
                </p>
                <p
                  className={`font-medium ${
                    new Date(codeDetails.expires_at) < new Date()
                      ? "text-red-400"
                      : ""
                  }`}
                >
                  {formatDate(codeDetails.expires_at)}
                </p>
              </div>

              <div>
                <p className="text-xs text-gray-400 font-medium mb-1">
                  Duration
                </p>
                <p className="font-medium">
                  {formatDate(codeDetails.duration) || "N/A"}
                </p>
              </div>

              <div>
                <p className="text-xs text-gray-400 font-medium mb-1">
                  Time Left
                </p>
                <p className="font-medium">
                  {codeDetails.expires_at &&
                  new Date(codeDetails.expires_at) > new Date()
                    ? `${Math.ceil(
                        (new Date(codeDetails.expires_at) - new Date()) /
                          (1000 * 60 * 60 * 24)
                      )} days`
                    : "Expired"}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-[#1e272e] p-5 rounded-lg">
            <h3 className="text-sm uppercase tracking-wider text-blue-400 font-semibold mb-3 border-b border-gray-700 pb-2">
              Usage Information
            </h3>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-xs text-gray-400 font-medium mb-1">
                  Usage Count
                </p>
                <p className="font-medium text-lg">
                  {codeDetails.used_count} / {codeDetails.max_uses}
                </p>
                <div className="w-full bg-gray-700 rounded-full h-2 mt-2">
                  <div
                    className="bg-blue-500 h-2 rounded-full"
                    style={{
                      width: `${
                        (codeDetails.used_count / codeDetails.max_uses) * 100
                      }%`,
                    }}
                  ></div>
                </div>
              </div>

              <div className="col-span-2">
                <p className="text-xs text-gray-400 font-medium mb-1">
                  Remaining Redemptions
                </p>
                <p className="font-medium text-lg">
                  {Math.max(0, codeDetails.max_uses - codeDetails.used_count)}
                  <span className="text-sm text-gray-400 ml-2">
                    (
                    {Math.round(
                      ((codeDetails.max_uses - codeDetails.used_count) /
                        codeDetails.max_uses) *
                        100
                    )}
                    % remaining)
                  </span>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {codeDetails.users && codeDetails.users.length > 0 && (
        <div className="p-6 bg-[#2d3436] rounded-lg">
          <h2 className="text-xl font-bold mb-4 flex items-center">
            <svg
              className="w-5 h-5 mr-2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
              ></path>
            </svg>
            Redemption History
            <span className="ml-2 bg-blue-900/30 text-blue-300 px-2 py-1 rounded-md text-xs font-medium">
              {codeDetails.users.length} users
            </span>
          </h2>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-[#1e272e] text-left">
                <tr>
                  <th
                    scope="col"
                    className="px-6 py-3 font-medium text-xs uppercase tracking-wider text-gray-300"
                  >
                    Username
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 font-medium text-xs uppercase tracking-wider text-gray-300"
                  >
                    Email
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 font-medium text-xs uppercase tracking-wider text-gray-300"
                  >
                    Status
                  </th>
                  <th
                    scope="col"
                    className="px-6 py-3 font-medium text-xs uppercase tracking-wider text-gray-300"
                  >
                    Redemption Date
                  </th>
                </tr>
              </thead>
              <tbody>
                {codeDetails.users.map((user) => (
                  <tr
                    key={user.id}
                    className="border-b border-gray-700 hover:bg-[#1e272e]/50"
                  >
                    <td className="px-6 py-4">
                      <a
                        href={`/admin/users/${user.id}`}
                        className="text-blue-400 hover:text-blue-300 hover:underline font-medium"
                      >
                        {user.username}
                      </a>
                    </td>
                    <td className="px-6 py-4 text-sm">{user.email}</td>
                    <td className="px-6 py-4">
                      {user.is_admin ? (
                        <span className="bg-indigo-900/30 text-indigo-300 border border-indigo-700 px-2 py-1 rounded-md text-xs">
                          Admin
                        </span>
                      ) : (
                        <span className="bg-blue-900/30 text-blue-300 border border-blue-700 px-2 py-1 rounded-md text-xs">
                          User
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm">
                      {formatDate(user.created_at)}
                    </td>
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
