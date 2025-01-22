const BASE_URL = "http://localhost:8000/api";

// Utility function for API requests
async function apiRequest(endpoint, method = "GET", data = null) {
    const options = { method, headers: {} };

    if (method !== "GET" && data) {
        if (data instanceof FormData) {
            options.body = data; // For file uploads
        } else {
            options.body = JSON.stringify(data);
            options.headers["Content-Type"] = "application/json";
        }
    }

    try {
        const response = await fetch(`${BASE_URL}${endpoint}`, options);
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || `HTTP Error: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error("API Request Error:", error);
        return { error: error.message };
    }
}

// Add Property
async function addProperty(formData) {
    const response = await apiRequest("/properties", "POST", formData);
    return response;
}


// Delete Property
async function deleteProperty(id) {
    if (!/^[a-fA-F0-9]{24}$/.test(id)) {
        return { error: "Invalid Property ID format. Must be a 24-character hexadecimal string." };
    }
    return await apiRequest(`/properties/${id}`, "DELETE");
}

// Get All Properties
async function getAllProperties() {
    return await apiRequest("/properties", "GET");
}

// Handle form submission for adding a property
document.getElementById("add-property-form")?.addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = new FormData(e.target);
    formData.append("user_id", "12345"); // Example static user_id

    const response = await addProperty(formData);

    // Show the message based on the response
    const responseMsg = document.getElementById("response-msg");
    if (response.error) {
        responseMsg.textContent = `Error: ${response.error}`;
        responseMsg.style.color = "red";
    } else {
        responseMsg.textContent = "Property added successfully!";
        responseMsg.style.color = "green";
        e.target.reset(); // Reset the form
    }

    // Refresh the property list
    const properties = await getAllProperties();
    renderProperties(properties);
});

// Update Property
async function updateProperty(id, formData) {
    return await apiRequest(`/properties/${id}`, "PUT", formData);
}

// Handle Update Form Submission
document.getElementById("propertyUpdateForm")?.addEventListener("submit", async (e) => {
    e.preventDefault();

    const propertyId = document.getElementById("id").value;
    const formData = new FormData(e.target);

    const response = await updateProperty(propertyId, formData);
    const responseMsg = document.getElementById("response-msg");

    if (response.error) {
        responseMsg.textContent = `Error: ${response.error}`;
        responseMsg.style.color = "red";
    } else {
        responseMsg.textContent = "Property updated successfully!";
        responseMsg.style.color = "green";
        e.target.reset();
    }
});


// Utility to render properties
function renderProperties(properties) {
    const propertiesContainer = document.getElementById("properties-container");
    if (!propertiesContainer) return;

    propertiesContainer.innerHTML = ""; // Clear existing properties

    properties.forEach((property) => {
        const propertyItem = document.createElement("div");
        propertyItem.className = "property-item";
        propertyItem.innerHTML = `
            <h3>${property.title}</h3>
            <p><strong>Location:</strong> ${property.location}</p>
            <p><strong>Price:</strong> ${property.price}</p>
            <p><strong>Type:</strong> ${property.propertyType}</p>
            <p><strong>Description:</strong> ${property.description}</p>
            <button onclick="handleDeleteProperty('${property.id}')">Delete</button>
        `;
        propertiesContainer.appendChild(propertyItem);
    });
}

// Handle deleting a property
async function handleDeleteProperty(id) {
    const response = await deleteProperty(id);
    if (response.error) {
        alert(`Error deleting property: ${response.error}`);
    } else {
        alert("Property deleted successfully!");
        const updatedProperties = await getAllProperties();
        renderProperties(updatedProperties);
    }
}

// Initial data fetch
document.addEventListener("DOMContentLoaded", async () => {
    const properties = await getAllProperties();
    renderProperties(properties);
});


