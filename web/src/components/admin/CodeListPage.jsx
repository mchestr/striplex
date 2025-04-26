import React, { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";

function CodeListPage({ onViewCodeDetails }) {
  const [inviteCodes, setInviteCodes] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [codeToDelete, setCodeToDelete] = useState(null);
  const [deleteError, setDeleteError] = useState(null); // New state for delete errors
  const [newCodeDuration, setNewCodeDuration] = useState("unlimited");
  const [newCodeMaxUses, setNewCodeMaxUses] = useState(1);
  const [createError, setCreateError] = useState(null);
  const [createSuccess, setCreateSuccess] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [showAdvancedOptions, setShowAdvancedOptions] = useState(false);
  const [newCodeDurationDate, setNewCodeDurationDate] = useState(1);
  const [newCodeExpirationDate, setNewCodeExpirationDate] = useState("");
  const [expirationOption, setExpirationOption] = useState("never");
  const [durationOption, setDurationOption] = useState("never");
  const customCodeInputRef = useRef(null); // Ref for the custom code input
  const [copiedCode, setCopiedCode] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchInviteCodes();
  }, []);

  const fetchInviteCodes = async () => {
    setIsLoading(true);
    try {
      const response = await fetch("/api/v1/codes");
      if (!response.ok) {
        if (response.status === 401) {
          navigate("/");
          return;
        }
        throw new Error(`Failed to fetch invite codes: ${response.status}`);
      }
      const data = await response.json();
      setInviteCodes(data.invite_codes || []);
    } catch (error) {
      console.error("Error fetching invite codes:", error);
      setError(error.message);
    } finally {
      setIsLoading(false);
    }
  };

  const calculateDate = (option) => {
    const today = new Date();
    switch (option) {
      case "7days":
        return new Date(today.setDate(today.getDate() + 7)).toISOString();
      case "1month":
        return new Date(today.setMonth(today.getMonth() + 1)).toISOString();
      case "3months":
        return new Date(today.setMonth(today.getMonth() + 3)).toISOString();
      case "6months":
        return new Date(today.setMonth(today.getMonth() + 6)).toISOString();
      case "12months":
        return new Date(today.setMonth(today.getMonth() + 12)).toISOString();
      default:
        return "";
    }
  };

  const handleExpirationOptionChange = (option) => {
    setExpirationOption(option);
    if (option === "never") {
      setNewCodeExpirationDate("");
    } else {
      setNewCodeExpirationDate(calculateDate(option));
    }
  };

  const handleDurationOptionChange = (option) => {
    setDurationOption(option);
    if (option === "never") {
      setNewCodeDurationDate("");
    } else {
      setNewCodeDurationDate(calculateDate(option));
    }
  };

  const handleCreateCode = async (e) => {
    e.preventDefault();
    setCreateError(null);
    setCreateSuccess(false);

    try {
      const currentCustomCodeValue = customCodeInputRef.current
        ? customCodeInputRef.current.value
        : "";

      let payload = {
        max_uses: parseInt(newCodeMaxUses, 10),
      };

      if (currentCustomCodeValue.trim()) {
        payload.code = currentCustomCodeValue.trim();
      }

      const duration =
        durationOption === "custom"
          ? newCodeDurationDate
          : calculateDate(durationOption);
      const expiration =
        expirationOption === "custom"
          ? newCodeExpirationDate
          : calculateDate(expirationOption);
      if (duration) {
        payload.duration = duration;
      }
      if (expiration) {
        payload.expires_at = expiration;
      }
      const response = await fetch("/api/v1/codes", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`Failed to create invite code: ${response.status}`);
      }

      await fetchInviteCodes();
      setCreateSuccess(false);
      setShowCreateModal(false);
      setExpirationOption("never");
      setDurationOption("never");
      if (customCodeInputRef.current) {
        customCodeInputRef.current.value = "";
      }
    } catch (error) {
      console.error("Error creating invite code:", error);
      setCreateError(error.message);
    }
  };

  const handleDeleteCode = (codeId) => {
    setCodeToDelete(codeId);
    setDeleteError(null); // Clear any previous error
    setShowDeleteModal(true);
  };

  const confirmDeleteCode = async () => {
    try {
      const response = await fetch(`/api/v1/codes/${codeToDelete}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error(`Failed to delete invite code: ${response.status}`);
      }

      setInviteCodes(inviteCodes.filter((code) => code.id !== codeToDelete));
      setShowDeleteModal(false);
      setCodeToDelete(null);
    } catch (error) {
      console.error("Error deleting invite code:", error);
      // Set error in state instead of showing alert
      setDeleteError(error.message);
    }
  };

  const handleCopyLink = (code) => {
    const claimUrl = `${window.location.origin}/claim/${code}`;
    navigator.clipboard
      .writeText(claimUrl)
      .then(() => {
        setCopiedCode(code);
        setTimeout(() => setCopiedCode(null), 2000); // Reset after 2 seconds
      })
      .catch((err) => {
        console.error("Failed to copy link: ", err);
      });
  };

  const formatDate = (timestamp) => {
    if (!timestamp) return "N/A";
    const date = new Date(timestamp);
    return date.toLocaleDateString(); // This shows only year/month/day in the local format
  };

  const filteredCodes = inviteCodes.filter((code) =>
    code.code.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleViewCodeDetails = (codeId) => {
    onViewCodeDetails(codeId);
  };

  const CreateCodeModal = () => (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-[#2d3436] rounded-lg p-6 max-w-md w-full mx-auto">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-xl font-bold">Create New Invite Code</h3>
          <button
            onClick={() => {
              setShowCreateModal(false);
              if (customCodeInputRef.current) {
                customCodeInputRef.current.value = "";
              }
            }}
            className="text-gray-400 hover:text-white"
          >
            ✕
          </button>
        </div>

        {createError && (
          <div className="mb-4 p-3 bg-red-500/20 border border-red-500/50 rounded-lg text-red-200">
            {createError}
          </div>
        )}

        {createSuccess && (
          <div className="mb-4 p-3 bg-green-500/20 border border-green-500/50 rounded-lg text-green-200">
            Invite code created successfully!
          </div>
        )}

        <form onSubmit={handleCreateCode}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-300 mb-1">
              Expiration
            </label>
            <select
              value={expirationOption}
              onChange={(e) => handleExpirationOptionChange(e.target.value)}
              className="w-full p-2.5 bg-[#3a4149] border border-gray-600 rounded-lg text-white focus:ring-blue-500 focus:border-blue-500"
              required
            >
              <option value="never">Unlimited (Never Expires)</option>
              <option value="7days">7 Days</option>
              <option value="1month">1 Month</option>
              <option value="3months">3 Months</option>
              <option value="6months">6 Months</option>
              <option value="12months">12 Months</option>
              <option value="custom">Custom Date</option>
            </select>

            {expirationOption === "custom" && (
              <input
                type="datetime-local"
                value={newCodeExpirationDate}
                onChange={(e) => setNewCodeExpirationDate(e.target.value)}
                min={new Date().toISOString().slice(0, 16)}
                className="w-full mt-2 p-2.5 bg-[#3a4149] border border-gray-600 rounded-lg text-white focus:ring-blue-500 focus:border-blue-500"
              />
            )}
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-300 mb-1">
              Max Uses
            </label>
            <input
              type="number"
              value={newCodeMaxUses}
              onChange={(e) => setNewCodeMaxUses(e.target.value)}
              min="1"
              max="100"
              className="w-full p-2.5 bg-[#3a4149] border border-gray-600 rounded-lg text-white focus:ring-blue-500 focus:border-blue-500"
              required
            />
          </div>

          <div className="mb-4">
            <button
              type="button"
              onClick={() => setShowAdvancedOptions(!showAdvancedOptions)}
              className="flex items-center text-sm text-gray-300 hover:text-white"
            >
              <svg
                className={`w-4 h-4 mr-1 transition-transform ${
                  showAdvancedOptions ? "rotate-180" : ""
                }`}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M19 9l-7 7-7-7"
                ></path>
              </svg>
              Advanced Options
            </button>

            {showAdvancedOptions && (
              <div className="mt-3 p-3 bg-[#1e272e]/30 border border-gray-700 rounded-lg">
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-300 mb-1">
                    Custom Code Value
                  </label>
                  <input
                    type="text"
                    defaultValue=""
                    placeholder="Leave empty for auto-generated code"
                    className="w-full p-2.5 bg-[#3a4149] border border-gray-600 rounded-lg text-white focus:ring-blue-500 focus:border-blue-500"
                    ref={customCodeInputRef}
                  />
                  <p className="mt-1 text-xs text-gray-400">
                    Specify a custom code or leave blank to auto-generate one
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-1">
                    Duration
                  </label>
                  <select
                    value={durationOption}
                    onChange={(e) => handleDurationOptionChange(e.target.value)}
                    className="w-full p-2.5 bg-[#3a4149] border border-gray-600 rounded-lg text-white focus:ring-blue-500 focus:border-blue-500"
                    required
                  >
                    <option value="never">Unlimited</option>
                    <option value="7days">7 Days</option>
                    <option value="1month">1 Month</option>
                    <option value="3months">3 Months</option>
                    <option value="6months">6 Months</option>
                    <option value="12months">12 Months</option>
                    <option value="custom">Custom Date</option>
                  </select>

                  {durationOption === "custom" && (
                    <input
                      type="datetime-local"
                      value={newCodeExpirationDate}
                      onChange={(e) => setNewCodeDuration(e.target.value)}
                      min={new Date().toISOString().slice(0, 16)}
                      className="w-full mt-2 p-2.5 bg-[#3a4149] border border-gray-600 rounded-lg text-white focus:ring-blue-500 focus:border-blue-500"
                    />
                  )}
                </div>
              </div>
            )}
          </div>

          <div className="flex justify-end">
            <button
              type="button"
              onClick={() => {
                setShowCreateModal(false);
                if (customCodeInputRef.current) {
                  customCodeInputRef.current.value = "";
                }
              }}
              className="mr-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-[#4b6bfb] hover:bg-[#3557fa] text-white rounded-lg"
            >
              Create
            </button>
          </div>
        </form>
      </div>
    </div>
  );

  const DeleteCodeModal = () => (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-[#2d3436] rounded-lg p-6 max-w-md w-full mx-auto">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-xl font-bold">Delete Invite Code</h3>
          <button
            onClick={() => {
              setShowDeleteModal(false);
              setDeleteError(null);
            }}
            className="text-gray-400 hover:text-white"
          >
            ✕
          </button>
        </div>

        {deleteError && (
          <div className="mb-4 p-3 bg-red-500/20 border border-red-500/50 rounded-lg text-red-200">
            {deleteError}
          </div>
        )}

        <p className="mb-6 text-gray-300">
          Are you sure you want to delete this invite code? This action cannot
          be undone.
        </p>

        <div className="flex justify-end">
          <button
            onClick={() => {
              setShowDeleteModal(false);
              setDeleteError(null);
            }}
            className="mr-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg"
          >
            Cancel
          </button>
          <button
            onClick={confirmDeleteCode}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-full">
        <div className="text-xl">Loading invite codes...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col justify-center items-center h-full">
        <div className="text-xl text-red-500">{error}</div>
      </div>
    );
  }

  return (
    <>
      {showCreateModal && <CreateCodeModal />}
      {showDeleteModal && <DeleteCodeModal />}

      <div className="mb-6">
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center">
          <h1 className="text-3xl font-bold mb-4 md:mb-0">Invite Codes</h1>

          <div className="flex flex-col sm:flex-row w-full md:w-auto gap-3">
            <div className="relative">
              <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                <svg
                  className="w-4 h-4 text-gray-400"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 20 20"
                >
                  <path
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"
                  />
                </svg>
              </div>
              <input
                type="search"
                placeholder="Search codes"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="bg-[#2d3436] border border-gray-600 text-white text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 p-2.5"
              />
            </div>

            <button
              onClick={() => setShowCreateModal(true)}
              className="px-4 py-2.5 bg-[#4b6bfb] hover:bg-[#3557fa] text-white rounded-lg shadow-md hover:shadow-lg transition font-medium flex items-center justify-center"
            >
              <svg
                className="w-4 h-4 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                ></path>
              </svg>
              Create Code
            </button>
          </div>
        </div>
      </div>

      <div className="bg-[#2d3436] shadow-lg rounded-lg overflow-hidden">
        <div className="overflow-x-auto">
          {filteredCodes.length === 0 ? (
            <div className="p-8 text-center">
              <p className="text-lg text-gray-400">
                {inviteCodes.length === 0
                  ? "No invite codes have been created yet."
                  : "No matching invite codes found."}
              </p>
              {inviteCodes.length === 0 && (
                <button
                  onClick={() => setShowCreateModal(true)}
                  className="mt-4 px-4 py-2 bg-[#4b6bfb] hover:bg-[#3557fa] text-white rounded-lg"
                >
                  Create your first code
                </button>
              )}
            </div>
          ) : (
            <table className="w-full">
              <thead className="bg-[#1e272e]/50 text-left">
                <tr className="border-b border-gray-700">
                  <th scope="col" className="px-6 py-4 font-medium">
                    Code
                  </th>
                  <th scope="col" className="px-6 py-4 font-medium">
                    Created
                  </th>
                  <th scope="col" className="px-6 py-4 font-medium">
                    Expires
                  </th>
                  <th scope="col" className="px-6 py-4 font-medium">
                    Duration
                  </th>
                  <th scope="col" className="px-6 py-4 font-medium">
                    Entitlement
                  </th>
                  <th scope="col" className="px-6 py-4 font-medium">
                    Uses
                  </th>
                  <th scope="col" className="px-6 py-4 font-medium">
                    Status
                  </th>
                  <th scope="col" className="px-6 py-4 font-medium">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {filteredCodes.map((code) => (
                  <tr
                    key={code.id}
                    className="border-b border-gray-700 hover:bg-[#2a3138] transition-colors"
                  >
                    <td className="px-6 py-4 font-mono">
                      <button
                        onClick={() => handleViewCodeDetails(code.id)}
                        className="text-blue-400 hover:text-blue-300 hover:underline text-left font-mono"
                      >
                        {code.code}
                      </button>
                    </td>
                    <td className="px-6 py-4">{formatDate(code.created_at)}</td>
                    <td className="px-6 py-4">{formatDate(code.expires_at)}</td>
                    <td className="px-6 py-4">{formatDate(code.duration)}</td>
                    <td className="px-6 py-4">
                      {code.entitlement_name || "N/A"}
                    </td>
                    <td className="px-6 py-4">
                      {code.used_count}/{code.max_uses}
                    </td>
                    <td className="px-6 py-4">
                      <span
                        className={`px-2 py-1 rounded-full text-xs ${
                          code.is_disabled
                            ? "bg-gray-800 text-gray-300"
                            : code.used_count >= code.max_uses
                            ? "bg-gray-800 text-gray-300"
                            : new Date(code.expires_at) < new Date()
                            ? "bg-red-900 text-red-300"
                            : "bg-green-900 text-green-300"
                        }`}
                      >
                        {code.is_disabled
                          ? "Disabled"
                          : code.used_count >= code.max_uses
                          ? "Used"
                          : new Date(code.expires_at) < new Date()
                          ? "Expired"
                          : "Active"}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex space-x-2">
                        {!code.is_disabled &&
                          code.used_count < code.max_uses &&
                          (!code.expires_at ||
                            new Date(code.expires_at) > new Date()) && (
                            <button
                              onClick={() => handleCopyLink(code.code)}
                              className="px-3 py-1 bg-blue-800 hover:bg-blue-700 text-blue-100 rounded-md text-sm flex items-center"
                              title="Copy shareable link"
                            >
                              {copiedCode === code.code ? (
                                <>
                                  <svg
                                    className="w-4 h-4 mr-1"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                    xmlns="http://www.w3.org/2000/svg"
                                  >
                                    <path
                                      fillRule="evenodd"
                                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                      clipRule="evenodd"
                                    ></path>
                                  </svg>
                                  Copied
                                </>
                              ) : (
                                <>
                                  <svg
                                    className="w-4 h-4 mr-1"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                    xmlns="http://www.w3.org/2000/svg"
                                  >
                                    <path d="M7 9a2 2 0 012-2h6a2 2 0 012 2v6a2 2 0 01-2 2H9a2 2 0 01-2-2V9z"></path>
                                    <path d="M5 3a2 2 0 00-2 2v6a2 2 0 002 2V5h8a2 2 0 00-2-2H5z"></path>
                                  </svg>
                                  Share
                                </>
                              )}
                            </button>
                          )}
                        <button
                          onClick={() => handleDeleteCode(code.id)}
                          className="px-3 py-1 bg-red-800 hover:bg-red-700 text-red-100 rounded-md text-sm"
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </>
  );
}

export default CodeListPage;
