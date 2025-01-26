import React, { useState, useEffect } from "react";
import "./dashboard.css";

const Dashboard = () => {
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
    try {
      setLoading(true);
      const response = await fetch("http://localhost:8080/api/user-details"); // Replace with your API endpoint
      const data = await response.json();
      setUserDetails(data);
    } catch (error) {
      console.error("Error fetching user details:", error);
    } finally {
      setLoading(false);
    }
  };

  const updateUserDetails = async () => {
    try {
      setLoading(true);
      const response = await fetch("http://localhost:8080/api/update-user-details", {
        method: "PUT",
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
    try {
      const response = await fetch("http://localhost:8080/api/get-qr-code"); // Replace with your API endpoint
      const blob = await response.blob();
      const qrCodeUrl = URL.createObjectURL(blob);
      setQrCode(qrCodeUrl);
    } catch (error) {
      console.error("Error fetching QR code:", error);
    }
  };

  useEffect(() => {
    fetchUserDetails();
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
