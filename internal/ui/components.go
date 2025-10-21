package ui

const (
	workflowColumnTemplate = `{{define "workflowColumn"}}
	<div class="column-header">
		<h2>Workflows</h2>
		<button id="open-workflow-modal" class="primary-action" type="button" title="Create workflow" aria-label="Create workflow">+</button>
	</div>
	{{if .}}
	<ul class="item-list">
		{{range .}}
		<li class="list-item">
			<details>
				<summary>
					<span class="item-title">{{.DisplayName}}</span>
				</summary>
				<div class="item-details">
					<p><strong>ID:</strong> {{.Workflow.ID}}</p>
					{{if .Workflow.Title}}<p><strong>Title:</strong> {{.Workflow.Title}}</p>{{end}}
					{{if .Workflow.Description}}<p><strong>Description:</strong> {{.Workflow.Description}}</p>{{end}}
					{{if .Workflow.Message}}<p><strong>Message:</strong> {{.Workflow.Message}}</p>{{end}}
					{{if .StepNames}}
					<h3>Steps</h3>
					<ul class="detail-list">
						{{range .StepNames}}
						<li>{{.}}</li>
						{{end}}
					</ul>
					{{end}}
				</div>
			</details>
		</li>
		{{end}}
	</ul>
	{{else}}
	<p>No workflows available.</p>
	{{end}}
{{end}}`

	runColumnTemplate = `{{define "runColumn"}}
	<h2>Runs</h2>
	{{if .}}
	<ul class="item-list">
		{{range .}}
		<li class="list-item">
			<details>
				<summary>
					<span class="item-title">{{.Name}}</span>
					<span class="item-state state {{.StateClass}}">{{.StateLabel}}</span>
				</summary>
				<div class="item-details">
					<p><strong>Workflow:</strong> {{.WorkflowName}}</p>
					{{if .Steps}}
					<h3>Steps</h3>
					<ul class="detail-list">
						{{range .Steps}}
						<li>
							<span>{{.Name}}</span>
							<span class="state {{.StatusClass}}">— {{.Status}}</span>
						</li>
						{{end}}
					</ul>
					{{end}}
				</div>
			</details>
		</li>
		{{end}}
	</ul>
	{{else}}
	<p>No runs found.</p>
	{{end}}
{{end}}`

	workflowModalTemplate = `{{define "workflowModal"}}
<div id="workflow-modal" class="modal-backdrop" role="dialog" aria-modal="true" aria-hidden="true">
	<div class="modal-card" role="document">
		<div class="modal-header">
			<h2>Create Workflow</h2>
			<button type="button" class="close-button" data-close-modal aria-label="Close create workflow form">×</button>
		</div>
		<div class="modal-body">
			<div id="workflow-form-error" class="error-banner" role="alert"></div>
			<form id="workflow-form">
				<div class="form-row">
					<label for="workflow-id">Workflow ID</label>
					<input id="workflow-id" name="workflow-id" type="text" placeholder="my-workflow" required>
				</div>
				<div class="form-row">
					<label for="workflow-title">Title</label>
					<input id="workflow-title" name="workflow-title" type="text" placeholder="Human-friendly workflow title" required>
				</div>
				<div class="form-row">
					<label for="workflow-description">Description</label>
					<textarea id="workflow-description" name="workflow-description" placeholder="Explain what this workflow accomplishes"></textarea>
				</div>
				<div class="form-row">
					<label for="workflow-message">Message</label>
					<textarea id="workflow-message" name="workflow-message" placeholder="Optional run message shown to operators"></textarea>
				</div>
				<div class="form-row">
					<div class="step-header">
						<h3>Steps</h3>
						<button type="button" id="add-workflow-step" class="add-step-button">+ Add Step</button>
					</div>
					<div id="workflow-steps"></div>
				</div>
				<div class="form-actions">
					<button type="button" class="secondary-button" data-close-modal>Cancel</button>
					<button type="submit" id="workflow-submit" class="primary-button">Save Workflow</button>
				</div>
			</form>
		</div>
	</div>
</div>
<template id="workflow-step-template">
	<div class="workflow-step">
		<div class="step-header">
			<h4>Step <span class="step-number"></span></h4>
			<button type="button" class="text-button remove-step">Remove</button>
		</div>
		<div class="form-row">
			<label>Step Name</label>
			<input type="text" name="step-name" placeholder="identify-step" required>
		</div>
		<div class="form-row">
			<label>Description</label>
			<textarea name="step-description" placeholder="Optional description"></textarea>
		</div>
		<div class="form-row">
			<label>Handler</label>
			<select name="step-handler">
				<option value="tool" selected>tool</option>
				<option value="human">human</option>
			</select>
		</div>
		<div class="form-row">
			<label>Prompt</label>
			<textarea name="step-prompt" placeholder="Guidance for human or cognitive handlers"></textarea>
		</div>
		<div class="form-row">
			<label>Content</label>
			<textarea name="step-content" placeholder="Optional inline content"></textarea>
		</div>
		<div class="form-row">
			<label>Inputs (one per line or comma separated)</label>
			<textarea name="step-inputs" placeholder="input-a&#10;input-b"></textarea>
		</div>
		<div class="form-row">
			<label>Output</label>
			<input type="text" name="step-output" placeholder="result-key">
		</div>
	</div>
</template>
{{end}}`

	dashboardPageTemplate = `{{define "dashboard"}}
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>Composer Workflow Dashboard</title>
	<style>
		body {
			margin: 0;
			background: #0f0f0f;
			color: #ffffff;
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
			line-height: 1.5;
		}
		main {
			margin: 1.5rem auto;
			max-width: 960px;
		}
		h1, h2 {
			margin-top: 0;
		}
		.columns {
			display: flex;
			gap: 1.5rem;
			align-items: flex-start;
			flex-wrap: wrap;
		}
		.column {
			background: #1a1a1a;
			flex: 1 1 280px;
			padding: 0.75rem 1rem;
			border-radius: 6px;
			box-sizing: border-box;
			position: relative;
		}
		.column ul {
			margin: 0;
			padding: 0;
		}
		.column-header {
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 0.75rem;
			margin-bottom: 0.75rem;
		}
		.primary-action {
			background: #2ecc71;
			border: none;
			color: #0f0f0f;
			border-radius: 999px;
			width: 2rem;
			height: 2rem;
			display: inline-flex;
			align-items: center;
			justify-content: center;
			font-size: 1.5rem;
			line-height: 1;
			cursor: pointer;
			transition: transform 0.12s ease, box-shadow 0.12s ease;
			box-shadow: 0 4px 10px rgba(46, 204, 113, 0.3);
		}
		.primary-action:focus {
			outline: 2px solid rgba(46, 204, 113, 0.9);
			outline-offset: 2px;
		}
		.primary-action:hover {
			transform: translateY(-1px);
			box-shadow: 0 6px 14px rgba(46, 204, 113, 0.35);
		}
		.primary-action:active {
			transform: translateY(0);
			box-shadow: 0 2px 6px rgba(46, 204, 113, 0.25);
		}
		.state {
			font-weight: 600;
		}
		.state-ready {
			color: #6cb4ff;
		}
		.state-succeeded {
			color: #65d57c;
		}
		.state-failed {
			color: #ff6b6b;
		}
		.state-pending {
			color: #ffd966;
		}
		.state-unknown {
			color: #bbbbbb;
		}
		.item-list {
			list-style: none;
		}
		.list-item {
			margin: 0.5rem 0;
			background: #232323;
			border-radius: 6px;
			border: 1px solid #2d2d2d;
			overflow: hidden;
		}
		.list-item details {
			display: block;
		}
		.list-item summary {
			cursor: pointer;
			padding: 0.6rem 0.8rem;
			margin: 0;
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 0.75rem;
			font-weight: 600;
		}
		.list-item summary::-webkit-details-marker {
			display: none;
		}
		.item-title {
			flex: 1;
		}
		.item-details {
			padding: 0.6rem 0.8rem 0.8rem;
			font-size: 0.95rem;
			color: #dddddd;
		}
		.item-details p {
			margin: 0.2rem 0;
		}
		.item-details h3 {
			margin: 0.6rem 0 0.3rem;
			font-size: 1rem;
		}
		.detail-list {
			list-style: none;
			padding: 0;
			margin: 0.2rem 0 0;
		}
		.detail-list li {
			margin: 0.2rem 0;
		}
		.detail-list .state {
			font-weight: 500;
		}
		.modal-backdrop {
			position: fixed;
			inset: 0;
			background: rgba(0, 0, 0, 0.75);
			display: none;
			align-items: flex-start;
			justify-content: center;
			overflow-y: auto;
			padding: 4rem 1rem;
			z-index: 1000;
		}
		.modal-backdrop.is-visible {
			display: flex;
		}
		.modal-card {
			background: #1a1a1a;
			border-radius: 10px;
			border: 1px solid #2d2d2d;
			max-width: 640px;
			width: 100%;
			box-shadow: 0 20px 42px rgba(0, 0, 0, 0.45);
		}
		.modal-header {
			padding: 1.25rem 1.5rem 0.75rem;
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 1rem;
			border-bottom: 1px solid #2d2d2d;
		}
		.modal-body {
			padding: 1rem 1.5rem 1.5rem;
		}
		.close-button {
			background: transparent;
			border: none;
			color: #bbbbbb;
			font-size: 1.5rem;
			cursor: pointer;
		}
		.close-button:focus {
			outline: 2px solid rgba(255, 255, 255, 0.3);
			outline-offset: 2px;
		}
		form label {
			display: block;
			font-weight: 600;
			margin-bottom: 0.25rem;
		}
		form input[type="text"],
		form textarea,
		form select {
			width: 100%;
			background: #0f0f0f;
			border: 1px solid #3a3a3a;
			border-radius: 6px;
			padding: 0.45rem 0.6rem;
			color: #ffffff;
			font-size: 0.95rem;
			box-sizing: border-box;
		}
		form textarea {
			min-height: 72px;
			resize: vertical;
		}
		.form-row {
			margin-bottom: 1rem;
		}
		.form-actions {
			margin-top: 1.5rem;
			display: flex;
			justify-content: flex-end;
			gap: 0.75rem;
		}
		.secondary-button {
			background: transparent;
			border: 1px solid #3a3a3a;
			color: #ffffff;
			padding: 0.45rem 0.9rem;
			border-radius: 6px;
			cursor: pointer;
		}
		.secondary-button:hover {
			background: rgba(255, 255, 255, 0.05);
		}
		.primary-button {
			background: #2ecc71;
			border: none;
			color: #0f0f0f;
			padding: 0.55rem 1.2rem;
			border-radius: 6px;
			cursor: pointer;
			font-weight: 600;
		}
		.primary-button:disabled {
			opacity: 0.6;
			cursor: not-allowed;
		}
		.workflow-step {
			border: 1px solid #2d2d2d;
			border-radius: 8px;
			padding: 0.9rem;
			background: #212121;
			margin-bottom: 1rem;
		}
		.step-header {
			display: flex;
			align-items: center;
			justify-content: space-between;
			margin-bottom: 0.6rem;
		}
		.text-button {
			background: transparent;
			border: none;
			color: #ff6b6b;
			cursor: pointer;
			font-weight: 600;
		}
		.add-step-button {
			margin-top: 0.5rem;
			display: inline-flex;
			align-items: center;
			gap: 0.4rem;
			background: transparent;
			border: 1px solid #2ecc71;
			color: #2ecc71;
			border-radius: 6px;
			padding: 0.4rem 0.75rem;
			cursor: pointer;
		}
		.error-banner {
			background: rgba(255, 107, 107, 0.12);
			border: 1px solid rgba(255, 107, 107, 0.5);
			color: #ffaaaa;
			border-radius: 6px;
			padding: 0.5rem 0.75rem;
			margin-bottom: 0.75rem;
			display: none;
		}
		.error-banner.is-visible {
			display: block;
		}
	</style>
</head>
<body>
	<main>
		<h1>Workflow Dashboard</h1>
		<div class="columns">
			<section class="column">
				{{template "workflowColumn" .Workflows}}
			</section>
			<section class="column">
				{{template "runColumn" .Runs}}
			</section>
		</div>
	</main>
	{{template "workflowModal" .}}
	<script>
	(function () {
		const modal = document.getElementById("workflow-modal");
		const form = document.getElementById("workflow-form");
		const openButton = document.getElementById("open-workflow-modal");
		const closeButtons = modal ? modal.querySelectorAll("[data-close-modal]") : [];
		const stepsContainer = document.getElementById("workflow-steps");
		const addStepButton = document.getElementById("add-workflow-step");
		const stepTemplate = document.getElementById("workflow-step-template");
		const errorBanner = document.getElementById("workflow-form-error");
		const submitButton = document.getElementById("workflow-submit");

		if (!modal || !form || !openButton || !stepsContainer || !addStepButton || !stepTemplate) {
			return;
		}

		const openModal = () => {
			resetForm();
			modal.classList.add("is-visible");
			modal.setAttribute("aria-hidden", "false");
			form.querySelector("input[name='workflow-id']").focus();
		};

		const closeModal = () => {
			modal.classList.remove("is-visible");
			modal.setAttribute("aria-hidden", "true");
		};

		const clearError = () => {
			errorBanner.textContent = "";
			errorBanner.classList.remove("is-visible");
		};

		const showError = (message) => {
			errorBanner.textContent = message;
			errorBanner.classList.add("is-visible");
		};

		const updateStepNumbers = () => {
			const stepElements = stepsContainer.querySelectorAll(".workflow-step");
			stepElements.forEach((el, index) => {
				const label = el.querySelector(".step-number");
				if (label) {
					label.textContent = index + 1;
				}
				const removeButton = el.querySelector(".remove-step");
				if (removeButton) {
					removeButton.disabled = stepElements.length === 1;
				}
			});
		};

		const addStep = (initialValues) => {
			const fragment = stepTemplate.content.cloneNode(true);
			const stepElement = fragment.querySelector(".workflow-step");
			const fields = {
				name: stepElement.querySelector("input[name='step-name']"),
				description: stepElement.querySelector("textarea[name='step-description']"),
				handler: stepElement.querySelector("select[name='step-handler']"),
				prompt: stepElement.querySelector("textarea[name='step-prompt']"),
				content: stepElement.querySelector("textarea[name='step-content']"),
				inputs: stepElement.querySelector("textarea[name='step-inputs']"),
				output: stepElement.querySelector("input[name='step-output']"),
			};

			if (initialValues) {
				if (fields.name) fields.name.value = initialValues.name || "";
				if (fields.description) fields.description.value = initialValues.description || "";
				if (fields.handler && initialValues.handler) fields.handler.value = initialValues.handler;
				if (fields.prompt) fields.prompt.value = initialValues.prompt || "";
				if (fields.content) fields.content.value = initialValues.content || "";
				if (fields.inputs && Array.isArray(initialValues.inputs)) {
					fields.inputs.value = initialValues.inputs.join("\n");
				}
				if (fields.output) fields.output.value = initialValues.output || "";
			}

			const removeButton = stepElement.querySelector(".remove-step");
			if (removeButton) {
				removeButton.addEventListener("click", () => {
					stepElement.remove();
					updateStepNumbers();
				});
			}

			stepsContainer.appendChild(stepElement);
			updateStepNumbers();
		};

		const resetForm = () => {
			form.reset();
			stepsContainer.innerHTML = "";
			clearError();
			addStep();
		};

		const parseInputs = (value) => {
			return value
				.split(/\r?\n|,/)
				.map((entry) => entry.trim())
				.filter(Boolean);
		};

		const handleSubmit = async (event) => {
			event.preventDefault();
			clearError();

			const id = form.querySelector("input[name='workflow-id']").value.trim();
			const title = form.querySelector("input[name='workflow-title']").value.trim();
			const description = form.querySelector("textarea[name='workflow-description']").value.trim();
			const message = form.querySelector("textarea[name='workflow-message']").value.trim();

			if (!id) {
				showError("Workflow ID is required.");
				return;
			}
			if (!title) {
				showError("Title is required.");
				return;
			}

			const steps = [];
			const stepElements = stepsContainer.querySelectorAll(".workflow-step");
			stepElements.forEach((stepElement) => {
				const name = stepElement.querySelector("input[name='step-name']").value.trim();
				const descriptionValue = stepElement.querySelector("textarea[name='step-description']").value.trim();
				const handler = stepElement.querySelector("select[name='step-handler']").value.trim();
				const prompt = stepElement.querySelector("textarea[name='step-prompt']").value.trim();
				const content = stepElement.querySelector("textarea[name='step-content']").value.trim();
				const inputsValue = stepElement.querySelector("textarea[name='step-inputs']").value.trim();
				const output = stepElement.querySelector("input[name='step-output']").value.trim();

				if (!name) {
					return;
				}

				const stepPayload = { name: name };
				if (descriptionValue) stepPayload.description = descriptionValue;
				if (handler) stepPayload.handler = handler;
				if (prompt) stepPayload.prompt = prompt;
				if (content) stepPayload.content = content;
				if (inputsValue) stepPayload.inputs = parseInputs(inputsValue);
				if (output) stepPayload.output = output;

				steps.push(stepPayload);
			});

			if (steps.length === 0) {
				showError("At least one step with a name is required.");
				return;
			}

			const payload = {
				title: title,
				description: description,
				message: message,
				steps: steps,
			};

			submitButton.disabled = true;
			try {
				const response = await fetch("/api/workflow/" + encodeURIComponent(id), {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify(payload),
				});

				if (!response.ok) {
					let messageText = "Failed to save workflow.";
					try {
						const data = await response.json();
						if (data && data.error && data.error.message) {
							messageText = data.error.message;
						}
					} catch (_) {
						/* ignore parse errors */
					}
					showError(messageText);
					return;
				}

				closeModal();
				window.location.reload();
			} catch (err) {
				showError("Network error while saving workflow.");
			} finally {
				submitButton.disabled = false;
			}
		};

		openButton.addEventListener("click", openModal);
		addStepButton.addEventListener("click", () => addStep());
		form.addEventListener("submit", handleSubmit);

		closeButtons.forEach((button) => {
			button.addEventListener("click", closeModal);
		});

		modal.addEventListener("click", (event) => {
			if (event.target === modal) {
				closeModal();
			}
		});

		document.addEventListener("keydown", (event) => {
			if (event.key === "Escape" && modal.classList.contains("is-visible")) {
				closeModal();
			}
		});
	})();
	</script>
</body>
</html>
{{end}}`
)
