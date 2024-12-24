document.addEventListener('DOMContentLoaded', function () {
    const modalTitle = document.getElementById('modal-title');
    const modalMessage = document.getElementById('modal-message');
    const resultLink = document.getElementById('result-link');
    const modalContainer = document.getElementById('modal-container');
    const closeModalBtn = document.getElementById('close-modal');
    const testButton = document.getElementById('test-button');
    const generateButton = document.getElementById('generate-button');

    // Ensure modal starts hidden
    modalContainer.classList.add('hidden');

    // Attach event listeners to both buttons with different endpoints
    testButton.addEventListener('click', () => handleButtonClick('/test'));
    generateButton.addEventListener('click', () => handleButtonClick('/generate'));

    async function handleButtonClick(endpoint) {
        // Show the modal with an initial message
        showModal("Processing your request...", "Connecting to server...");

        const query = document.getElementById('text-box').value; // Get user input
        // if (!query) {
        //     modalMessage.textContent = "Please enter text in the box.";
        //     return;
        // }

        try {
            // Make the POST request to the specified endpoint
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ query })
            });

            if (!response.ok) throw new Error(`Failed to process request at ${endpoint}`);

            const responseData = await response.json();
            const { url, message } = responseData;

            // Update the modal with the result
            modalTitle.textContent = message || "Success!";
            modalMessage.textContent = "Your page is ready. Click below to view it:";
            resultLink.innerHTML = `<a href="${url}" target="_blank">${url}</a>`;

        } catch (error) {
            console.error(error);

            // Show error message in the modal
            modalTitle.textContent = "Error";
            modalMessage.textContent = "Failed to process your request. Please try again later.";
            resultLink.innerHTML = ""; // Clear any previous link
        }
    }

    // Helper function to show the modal
    function showModal(title, msg) {
        modalTitle.textContent = title;
        modalMessage.innerHTML = `<div class="loader"></div>${msg}`;
        resultLink.innerHTML = ""; // Clear any old content
        modalContainer.classList.remove('hidden'); // Show the modal
        modalContainer.classList.add('show'); // Ensure visibility
    }

    // Helper function to hide the modal
    function hideModal() {
        modalContainer.classList.add('hidden'); // Add "hidden" class to hide the modal
        modalContainer.classList.remove('show'); // Remove "show" class
    }

    // Close the modal when the "Close" button is clicked
    closeModalBtn.addEventListener('click', hideModal);
});
