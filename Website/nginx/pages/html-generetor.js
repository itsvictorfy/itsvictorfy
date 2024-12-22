document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('generate-form');
    const modalTitle = document.getElementById('modal-title');
    const modalMessage = document.getElementById('modal-message');
    const resultLink = document.getElementById('result-link');
    const modalContainer = document.getElementById('modal-container');
    const closeModalBtn = document.getElementById('close-modal');

    // Ensure modal starts hidden
    modalContainer.classList.add('hidden');

    form.addEventListener('submit', async function (event) {
        event.preventDefault(); // Prevent default form submission

        // Show the modal and set initial text
        showModal("Generating your page...", "Connecting to server...");

        const query = document.getElementById('text-box').value;
        let intervalId;
        try {
            // Start polling for progress
            const pollProgress = async (id) => {
                try {
                    const progressRes = await fetch(`/progress/${id}`);
                    if (progressRes.ok) {
                        const progressText = await progressRes.text();
                        modalMessage.textContent = progressText;
                    }
                } catch (err) {
                    console.error("Error fetching progress:", err);
                }
            };

            // Start the generation process
            const startRes = await fetch('/generate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ query })
            });

            if (!startRes.ok) throw new Error("Failed to start generation");

            const startData = await startRes.json();
            const { id, url, message } = startData;

            // Stop polling when the generation is complete
            clearInterval(intervalId);

            // Update the modal with the result
            modalTitle.textContent = message;
            modalMessage.textContent = "Your page is ready. Click below to view it:";
            resultLink.innerHTML = `<a href="${url}" target="_blank">${url}</a>`;

        } catch (err) {
            console.error(err);
            modalMessage.textContent = "Failed to generate page. Please try again later.";
        } finally {
            clearInterval(intervalId); // Ensure polling stops
        }
    });

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
