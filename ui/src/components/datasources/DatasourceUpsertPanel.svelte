<script>
    import { createEventDispatcher } from "svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import { setErrors } from "@/stores/errors";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast } from "@/stores/toasts";
    import Field from "@/components/base/Field.svelte";
    import Toggler from "@/components/base/Toggler.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Select from "@/components/base/Select.svelte";



    const dispatch = createEventDispatcher();
    const formId = "datasrouce_" + CommonHelper.randomString(5);

    let panel;

    let isSaving = false;
    let isConnectionTesting = false;
    let confirmClose = false; // prevent close recursion

    let datasource = {isNew:true};
    let name = "";
    let dsn = "";
    let type = "";

    let field = {options:{values:['mysql']},required:true}



    $: hasChanges =
        (datasource.isNew && datasource.name != "")||
        name != datasource.name ||
        dsn != datasource.dsn ||
        type!= datasource.type ;



    export function show(model) {
        load(model);
        
        confirmClose = true;

        return panel?.show();
    }

    export function hide() {
        return panel?.hide();
    }

    function load(model) {
        setErrors({}); // reset errors
        datasource = model?.clone ? model.clone() : {isNew:true};
        reset(); // reset form
    }

    function reset() {
        name = datasource?.name || "";
        dsn = datasource?.dsn || "";
        type = datasource?.type ||  "";
    }

    function save() {
        if (isSaving || !hasChanges) {
            return;
        }

        isSaving = true;

        const data = { name,dsn,type };

        let request;
        if (datasource.isNew) {
            request = ApiClient.records.create("datasources",data)
        } else {
            request = ApiClient.records.update("datasources",datasource.id, data);
        }

        request
            .then(async (result) => {
                confirmClose = false;
                hide();
                addSuccessToast(datasource.isNew ? "Successfully created datasource." : "Successfully updated datasource.");
                dispatch("save", result);

            })
            .catch((err) => {
                ApiClient.errorResponseHandler(err);
            })
            .finally(() => {
                isSaving = false;
            });
    }

    function deleteConfirm() {
        if (!datasource?.id) {
            return; // nothing to delete
        }

        confirm(`Do you really want to delete the selected datasource?`, () => {
            return ApiClient.records.delete("datasources",datasource.id)
                .then(() => {
                    confirmClose = false;
                    hide();
                    addSuccessToast("Successfully deleted datasource.");
                    dispatch("delete", datasource);
                })
                .catch((err) => {
                    ApiClient.errorResponseHandler(err);
                });
        });
    }

    function connectionTest() {
        if(isConnectionTesting){
            return
        }
        isConnectionTesting = true;
        const data = {name,dsn,type };

        let request  = ApiClient.send("/api/settings/test/connectdb",{
            'method': 'POST',
            'body':   data,
        });

        request
            .then(async (result) => {
                confirmClose = false;
                addSuccessToast("Successfully connect database." );

            })
            .catch((err) => {
                ApiClient.errorResponseHandler(err);
            })
            .finally(() => {
                isConnectionTesting = false;
            });

    }
</script>

<OverlayPanel
    bind:this={panel}
    popup
    class="datasource-panel"
    beforeHide={() => {
        if (hasChanges && confirmClose) {
            confirm("You have unsaved changes. Do you really want to close the panel?", () => {
                confirmClose = false;
                hide();
            });
            return false;
        }
        return true;
    }}
    on:hide
    on:show
>
    <svelte:fragment slot="header">
        <h4>
            {datasource.isNew ? "New datasource" : "Edit datasource"}
        </h4>
    </svelte:fragment>

    <form id={formId} class="grid" autocomplete="off" on:submit|preventDefault={save}>
        {#if !datasource.isNew}
            <Field class="form-field disabled" name="id" let:uniqueId>
                <label for={uniqueId}>
                    <i class={CommonHelper.getFieldTypeIcon("primary")} />
                    <span class="txt">ID</span>
                </label>
                <input type="text" id={uniqueId} value={datasource.id} disabled />
            </Field>
        {/if}


        <Field class="form-field required" name="name" let:uniqueId>
            <label for={uniqueId}>
                <i class={CommonHelper.getFieldTypeIcon("text")} />
                <span class="txt">Name</span>
            </label>
            <input type="text" autocomplete="off" id={uniqueId} required bind:value={name} />
        </Field>

        <Field class="form-field required" name="dsn" let:uniqueId>
            <label for={uniqueId}>
                <i class={CommonHelper.getFieldTypeIcon("text")} />
                <span class="txt">Dsn</span>
            </label>
            <input type="text" autocomplete="off" id={uniqueId} required bind:value={dsn} />
        </Field>


        <Field class="form-field required" name="type" let:uniqueId>
            <label for={uniqueId}>
                <i class={CommonHelper.getFieldTypeIcon("select")} />
                <span class="txt">Database Type</span>
            </label>
            <Select
                    id={uniqueId}
                    toggle={false}
                    items={field.options?.values}
                    searchable={field.options?.values > 5}
                    bind:selected={type}
            />
        </Field>


    </form>

    <svelte:fragment slot="footer">
        {#if !datasource.isNew}
            <button type="button" class="btn btn-sm btn-circle btn-secondary">
                <!-- empty span for alignment -->
                <span />
                <i class="ri-more-line" />
                <Toggler class="dropdown dropdown-upside dropdown-left dropdown-nowrap">
                    <button type="button" class="dropdown-item" on:click={() => deleteConfirm()}>
                        <i class="ri-delete-bin-7-line" />
                        <span class="txt">Delete</span>
                    </button>
                </Toggler>
            </button>
            <div class="flex-fill" />
        {/if}

        <button type="button" class="btn btn-secondary" disabled={isSaving} on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button
                type="button"
                class="btn btn-expanded"
                disabled={dsn=="" || isSaving || isConnectionTesting}
                on:click={()=>connectionTest()}
        >
            <span class="txt">Test connection</span>
        </button>
        <button
            type="submit"
            form={formId}
            class="btn btn-expanded"
            class:btn-loading={isSaving}
            disabled={!hasChanges || isSaving}
        >
            <span class="txt">{datasource.isNew ? "Create" : "Save changes"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>
