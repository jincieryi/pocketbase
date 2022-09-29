<script>
    import SortHeader from "@/components/base/SortHeader.svelte";
    import RecordFieldCell from "@/components/records/RecordFieldCell.svelte";
    import {Collection, Record} from "pocketbase";
    import {createEventDispatcher, tick} from "svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import {setErrors} from "@/stores/errors";
    import {confirm} from "@/stores/confirmation";
    import {addSuccessToast} from "@/stores/toasts";
    import {activeCollection, addCollection} from "@/stores/collections";
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import DataSourceSelect from "@/components/collections/DatasourceSelect.svelte"
    import SqlConfirmPanel from "@/components/collections/SqlConfirmPanel.svelte";




    let collection = new Collection();
    collection.exp=new Record();

    let selectedDatasource;

    $:{
        collection.exp.did = selectedDatasource?.id
    }


    let sqlConfirmPanel;



    let excuteErrorMsg = "";
    let isSearching = false;
    let searchResult = {};

    $: fields = searchResult?.schema||[];
    $: records = searchResult?.records||[];

    function search() {
        excuteErrorMsg=""

        if(!collection?.exp?.rawSql || !collection?.exp?.did){
            excuteErrorMsg = "Datasource or sql is invalid";
            return;
        }else{
            excuteErrorMsg = "";
        }

        isSearching = true;

        let request  = ApiClient.send("/api/collections/mountdb/excute-sql",{
            'method': 'POST',
            'body':   {rawSql:collection?.exp?.rawSql,did:collection?.exp?.did},
        });

         return request
            .then(async (result) => {
                isSearching = true;
                searchResult = result;
            })
            .catch((err) => {
                excuteErrorMsg = err
                isSearching = false;
                ApiClient.errorResponseHandler(err);
            })
            .finally(() => {
                isSearching = false;
            });

    }


    function asCollection() {
        search().then(
            (result)=>{
                if(excuteErrorMsg){
                   return
                }
                let c = collection?.clone()
                c.schema = searchResult.schema
                sqlConfirmPanel?.show(c)
                console.log(c)
            }
        );
    }



    const TAB_FIELDS = "fields";
    const dispatch = createEventDispatcher();

    let sqlCollectionPanel;


    let original = null;

    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let activeTab = TAB_FIELDS;
    let initialFormHash = calculateFormHash(collection);



    $: isSystemUpdate = !collection.isNew && collection.system;

    $: hasChanges = initialFormHash != calculateFormHash(collection);

    $: canSave = collection.isNew || hasChanges;

    export function changeTab(newTab) {
        activeTab = newTab;
    }

    export function show(model) {
        load(model);

        confirmClose = true;

        changeTab(TAB_FIELDS);

        return sqlCollectionPanel?.show();
    }

    export function hide() {
        return sqlCollectionPanel?.hide();
    }

    async function load(model) {
        setErrors({}); // reset errors
        if (typeof model !== "undefined") {
            //TODO get exp by cid
            original = model;
            collection = model?.clone();

        } else {
            original = null;
            collection.exp = new Record();

        }

    }



    function save() {
        if (isSaving) {
            return;
        }

        isSaving = true;

        const data = exportFormData();

        let request;
        if (collection.isNew) {
            request = ApiClient.collections.create(data);
        } else {
            request = ApiClient.collections.update(collection.id, data);
        }

        request
            .then((result) => {
                confirmClose = false;
                hide();
                addSuccessToast(
                    collection.isNew ? "Successfully created collection." : "Successfully updated collection."
                );
                addCollection(result);

                if (collection.isNew) {
                    $activeCollection = result;
                }

                dispatch("save", result);
            })
            .catch((err) => {
                ApiClient.errorResponseHandler(err);
            })
            .finally(() => {
                isSaving = false;
            });
    }

    function exportFormData() {
        const data = collection.export();
        data.schema = data.schema.slice(0);

        // remove deleted fields
        for (let i = data.schema.length - 1; i >= 0; i--) {
            const field = data.schema[i];
            if (field.toDelete) {
                data.schema.splice(i, 1);
            }
        }

        return data;
    }


    function calculateFormHash(m) {
        return JSON.stringify(m);
    }






</script>

<OverlayPanel
        bind:this={sqlCollectionPanel}
        class="overlay-panel-xl colored-header compact-header collection-panel"
        beforeHide={() => {
        if (!!collection?.exp?.rawSql && confirmClose) {
            confirm("You have unsaved changes. Do you really want to close the panel?", () => {
                confirmClose = false;
                searchResult={};
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
            {collection.isNew ? "New sqlcollection" : "Edit sqlcollection"}
        </h4>



        <form
                class="block"
                on:submit|preventDefault={search}
        >
            <div class="grid m-b-base">
                <div class="col-lg-6">
                    <Field
                            class="form-field required m-b-0 {isSystemUpdate ? 'disabled' : ''}"
                            name="name"
                            let:uniqueId
                    >
                        <label for={uniqueId}>Datasource</label>
                        <!-- svelte-ignore a11y-autofocus -->
                        <DataSourceSelect bind:selected={selectedDatasource}></DataSourceSelect>

                    </Field>
                </div>
            </div>

            <div class="grid m-b-base">
                <div class="col-lg-10">

                    <Field class="form-field required m-b-0 {isSystemUpdate ? 'disabled' : ''}"
                           name="name"
                           let:uniqueId>
                        <label for={uniqueId}>Sql Editor</label>

                        <textarea
                                id={uniqueId}
                                class="txt-mono"
                                spellcheck="false"
                                rows="8"
                                required
                                bind:value={collection.exp.rawSql}

                        />
                        {#if  !!excuteErrorMsg}
                            <div class="help-block help-block-error">{excuteErrorMsg}</div>
                        {/if}
                    </Field>
                </div>
                <div class="col-lg-2">
                    <button
                            type="button"
                            class="btn btn-block btn-success"
                            class:active={isSearching}
                            disabled={isSearching}
                            on:click={() => search()}
                    >
                        <span class="txt">Excute</span>
                    </button>
                    <br>
                </div>
            </div>

        </form>


    </svelte:fragment>

    <div class="table-wrapper">
        <table class="table" class:table-loading={isSearching}>
            <thead>
            <tr>
                {#each fields as field (field.name)}
                    <SortHeader
                            class="col-type-{field.type} col-field-{field.name}"
                            name={field.name}
                            disable="true"
                    >
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon(field.type)} />
                            <span class="txt">{field.name}</span>
                        </div>
                    </SortHeader>
                {/each}

                <th class="col-type-action min-width" />
            </tr>
            </thead>
            <tbody>
            {#each records as record }
                <tr tabindex="0" class="row-handle" >

                    {#each fields as field (field.name)}
                        <RecordFieldCell {record} {field} />
                    {/each}


                    <td class="col-type-action min-width">

                    </td>
                </tr>
            {:else}
                {#if isSearching}
                    <tr>
                        <td colspan="99" class="p-xs">
                            <span class="skeleton-loader" />
                        </td>
                    </tr>
                {:else}
                    <tr>
                        <td colspan="99" class="txt-center txt-hint p-xs">
                            <h6>No records found.</h6>
                        </td>
                    </tr>
                {/if}
            {/each}
            </tbody>
        </table>
    </div>

    <svelte:fragment slot="footer">
        <button type="button" class="btn btn-secondary"  on:click={() => hide()}>
            <span class="txt">Cancel</span>
        </button>
        <button
                type="button"
                class="btn btn-expanded"
                on:click={() => asCollection()}
        >
            <span class="txt">{collection.exp.isNew ? "As Collection" : "Save changes"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<SqlConfirmPanel bind:this={sqlConfirmPanel}/>
<style>
    .tabs-content {
        z-index: 3; /* autocomplete dropdown overlay fix */
    }

</style>
