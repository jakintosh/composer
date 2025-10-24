(function () {
  initWorkflowModal();
  initRunModal();
  initRunTickButtons();

  function initWorkflowModal() {
    const modal = document.getElementById("workflow-modal");
    const form = document.getElementById("workflow-form");
    const openButton = document.getElementById("open-workflow-modal");
    const stepsContainer = document.getElementById("workflow-steps");
    const addStepButton = document.getElementById("add-workflow-step");
    const stepTemplate = document.getElementById("workflow-step-template");
    const errorBanner = document.getElementById("workflow-form-error");
    const submitButton = document.getElementById("workflow-submit");

    if (!modal || !form || !openButton || !stepsContainer || !addStepButton || !stepTemplate || !errorBanner || !submitButton) {
      return;
    }

    const closeButtons = modal.querySelectorAll("[data-close-modal]");

    const openModal = () => {
      resetForm();
      modal.classList.add("is-visible");
      modal.setAttribute("aria-hidden", "false");
      const initialFocus = form.querySelector("input[name='workflow-id']");
      if (initialFocus) {
        initialFocus.focus();
      }
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
      if (!stepElement) {
        return;
      }

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
        if (fields.inputs) fields.inputs.value = Array.isArray(initialValues.inputs) ? initialValues.inputs.join("\n") : initialValues.inputs || "";
        if (fields.output) fields.output.value = initialValues.output || "";
      }

      const removeButton = stepElement.querySelector(".remove-step");
      if (removeButton) {
        removeButton.addEventListener("click", () => {
          stepElement.remove();
          updateStepNumbers();
        });
      }

      stepsContainer.appendChild(fragment);
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
            // ignore JSON parse errors
          }
          showError(messageText);
          return;
        }

        closeModal();
        window.location.reload();
      } catch (_err) {
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
  }

  function initRunModal() {
    const modal = document.getElementById("run-modal");
    const form = document.getElementById("run-form");
    const workflowIdInput = document.getElementById("run-workflow-id");
    const workflowNameInput = document.getElementById("run-workflow-name");
    const runIdInput = document.getElementById("run-id");
    const runNameInput = document.getElementById("run-name");
    const runIdPreview = document.getElementById("run-id-preview");
    const errorBanner = document.getElementById("run-form-error");
    const submitButton = document.getElementById("run-submit");
    const openButtons = document.querySelectorAll("[data-open-run-modal]");

    if (!modal || !form || !workflowIdInput || !workflowNameInput || !runIdInput || !runNameInput || !runIdPreview || !errorBanner || !submitButton || openButtons.length === 0) {
      return;
    }

    const closeButtons = modal.querySelectorAll("[data-close-modal]");

    const clearError = () => {
      errorBanner.textContent = "";
      errorBanner.classList.remove("is-visible");
    };

    const showError = (message) => {
      errorBanner.textContent = message;
      errorBanner.classList.add("is-visible");
    };

    const generateRunId = (raw) => {
      if (!raw) {
        return "";
      }
      const normalized = typeof raw.normalize === "function"
        ? raw.normalize("NFKD").replace(/[\u0300-\u036f]/g, "")
        : raw;
      return normalized
        .toLowerCase()
        .trim()
        .replace(/[^a-z0-9]+/g, "-")
        .replace(/^-+|-+$/g, "");
    };

    const updateRunId = () => {
      const displayName = runNameInput.value;
      const generated = generateRunId(displayName);
      runIdInput.value = generated;
      runIdPreview.textContent = generated || "—";
    };

    const resetForm = () => {
      form.reset();
      clearError();
      runIdPreview.textContent = "—";
      runIdInput.value = "";
    };

    const openModal = (workflowId, workflowLabel) => {
      resetForm();
      workflowIdInput.value = workflowId;
      workflowNameInput.value = workflowLabel || workflowId;
      workflowIdInput.value = workflowId;
      const baseLabel = workflowLabel || workflowId;
      if (baseLabel) {
        runNameInput.value = baseLabel + " Run";
      }
      updateRunId();
      modal.classList.add("is-visible");
      modal.setAttribute("aria-hidden", "false");
      runNameInput.focus();
      runNameInput.select();
    };

    const closeModal = () => {
      modal.classList.remove("is-visible");
      modal.setAttribute("aria-hidden", "true");
    };

    const handleSubmit = async (event) => {
      event.preventDefault();
      clearError();

      const workflowId = workflowIdInput.value.trim();
      const displayName = runNameInput.value.trim();
      const runId = runIdInput.value.trim();

      if (!workflowId) {
        showError("Workflow is required to start a run.");
        return;
      }
      if (!displayName) {
        showError("Display name is required.");
        return;
      }
      if (!runId) {
        showError("Display name must include at least one letter or number.");
        return;
      }

      const payload = {
        workflow_id: workflowId,
        name: displayName,
      };

      submitButton.disabled = true;
      try {
        const response = await fetch("/api/run/" + encodeURIComponent(runId), {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });

        if (!response.ok) {
          let messageText = "Failed to create run.";
          try {
            const data = await response.json();
            if (data && data.error && data.error.message) {
              messageText = data.error.message;
            }
          } catch (_) {
            // ignore JSON parse errors
          }
          showError(messageText);
          return;
        }

        closeModal();
        window.location.reload();
      } catch (_err) {
        showError("Network error while creating run.");
      } finally {
        submitButton.disabled = false;
      }
    };

    runNameInput.addEventListener("input", updateRunId);

    openButtons.forEach((button) => {
      button.addEventListener("click", (event) => {
        event.preventDefault();
        event.stopPropagation();
        const workflowId = button.getAttribute("data-workflow-id") || "";
        if (!workflowId) {
          return;
        }
        const workflowLabel = button.getAttribute("data-workflow-display") || workflowId;
        openModal(workflowId, workflowLabel);
      });
    });

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
  }

  function initRunTickButtons() {
    const buttons = document.querySelectorAll(".run-tick-button");
    if (buttons.length === 0) {
      return;
    }

    buttons.forEach((button) => {
      button.addEventListener("click", async (event) => {
        event.preventDefault();
        event.stopPropagation();

        const runId = button.getAttribute("data-run-id");
        if (!runId) {
          return;
        }

        const originalLabel = button.textContent;
        const runDisplay = button.getAttribute("data-run-display") || runId;

        button.classList.remove("has-error");
        button.removeAttribute("title");
        button.disabled = true;
        button.textContent = "Ticking...";

        try {
          const response = await fetch("/api/run/" + encodeURIComponent(runId) + "/tick", {
            method: "POST",
          });

          if (!response.ok) {
            let message = "Failed to tick " + runDisplay + ".";
            try {
              const data = await response.json();
              if (data && data.error && data.error.message) {
                message = data.error.message;
              }
            } catch (_ignore) {
              // Ignore JSON parsing errors on failure
            }
            throw new Error(message);
          }

          window.location.reload();
        } catch (error) {
          const message = error instanceof Error && error.message ? error.message : "Failed to tick " + runDisplay + ".";
          button.disabled = false;
          button.textContent = originalLabel;
          button.classList.add("has-error");
          button.setAttribute("title", message);
          console.error(message);
        }
      });
    });
  }
})();
