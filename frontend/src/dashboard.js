import React, { useState, useEffect } from "react";
import { useLocation } from "react-router-dom";
import "./dashboard.css";

const backendUrl = "http://localhost:8080/api/users";

const Dashboard = () => {
  const location = useLocation();
  const { email: userEmail } = location.state || {}; // Retrieve email from state
  const [userDetails, setUserDetails] = useState({
    email: "",
    phone: "",
    linkedin: "",
    github: "",
    leetcode: "",
  });
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(false);
  const [qrCode, setQrCode] = useState("");

  const fetchUserDetails = async () => {
    if (!userEmail) return; // Ensure email is available
    try {
      setLoading(true);
      const response = await fetch(backendUrl + "/details/" + userEmail); // Replace with your API endpoint
      const data = await response.json();
      setUserDetails(data);
    } catch (error) {
      console.error("Error fetching user details:", error);
    } finally {
      setLoading(false);
    }
  };

  const updateUserDetails = async () => {
    if (!userEmail) return; // Ensure email is available
    try {
      setLoading(true);
      const response = await fetch(backendUrl + "/details/" + userEmail, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(userDetails),
      });

      if (response.ok) {
        alert("Details updated successfully!");
        setIsEditing(false);
      } else {
        const errorData = await response.json();
        alert("Failed to update details: " + errorData.message);
      }
    } catch (error) {
      console.error("Error updating user details:", error);
    } finally {
      setLoading(false);
    }
  };

  const fetchQrCode = async () => {
    if (!userEmail) return; // Ensure email is available
    try {
      const response = await fetch(backendUrl + "/my-qr/" + userEmail); // Replace with your API endpoint
      const blob = await response.blob();
      const qrCodeUrl = URL.createObjectURL(blob);
      setQrCode(qrCodeUrl);
    } catch (error) {
      console.error("Error fetching QR code:", error);
    }
  };

  useEffect(() => {
    fetchUserDetails(); // Fetch user details on component mount
  }, []);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setUserDetails({ ...userDetails, [name]: value });
  };

  const handleCancel = () => {
    setIsEditing(false);
    fetchUserDetails(); // Reset the fields to the original values from the DB
  };

  return (
    <div className="dashboard">
      <br />
      <br />
      <br />
      <br />
      {loading && <p>Loading...</p>}
      {!loading && (
        <>
          <div className="user-details">
            <label>
              Email:
              <input
                type="text"
                name="email"
                value={userDetails.email}
                onChange={handleInputChange}
                disabled={!isEditing}
              />
            </label>
            <br />
            <label>
              Phone:
              <input
                type="text"
                name="phone"
                value={userDetails.phone}
                onChange={handleInputChange}
                disabled={!isEditing}
              />
            </label>
            <br />
            <label>
              LinkedIn:
              <input
                type="text"
                name="linkedin"
                value={userDetails.linkedin}
                onChange={handleInputChange}
                disabled={!isEditing}
              />
            </label>
            <br />
            <label>
              GitHub:
              <input
                type="text"
                name="github"
                value={userDetails.github}
                onChange={handleInputChange}
                disabled={!isEditing}
              />
            </label>
            <br />
            <label>
              LeetCode:
              <input
                type="text"
                name="leetcode"
                value={userDetails.leetcode}
                onChange={handleInputChange}
                disabled={!isEditing}
              />
            </label>
          </div>
          <div className="actions">
            {isEditing ? (
              <>
                <button onClick={updateUserDetails} disabled={loading}>
                  Save Changes
                </button>
                <button onClick={handleCancel} disabled={loading}>
                  Cancel
                </button>
              </>
            ) : (
              <button onClick={() => setIsEditing(true)}>Edit Details</button>
            )}
            {!isEditing && (
              <button onClick={fetchQrCode} disabled={loading}>
                My QR Code
              </button>
            )}
          </div>
          {qrCode && (
            <div className="qr-code">
              <h2>Your QR Code:</h2>
              <img src={qrCode} alt="QR Code" />
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default Dashboard;
