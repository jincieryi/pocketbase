<script>
    import SortHeader from "@/components/base/SortHeader.svelte";
    import RecordFieldCell from "@/components/records/RecordFieldCell.svelte";
    import {createEventDispatcher, tick} from "svelte";
    import CommonHelper from "@/utils/CommonHelper";
    import ApiClient from "@/utils/ApiClient";
    import {setErrors} from "@/stores/errors";
    import {resetConfirmation,confirm,confirmation} from "@/stores/confirmation";
    import {CollectionExp} from "@/stores/collections";
    import Field from "@/components/base/Field.svelte";
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import DataSourceSelect from "@/components/collections/DatasourceSelect.svelte"
    import SqlConfirmPanel from "@/components/collections/SqlConfirmPanel.svelte";



    let original = null;
    let collection = new CollectionExp();

    let selectedDatasource;

    $:{
        collection.did = selectedDatasource?.id
    }


    let sqlConfirmPanel;



    let excuteErrorMsg = "";
    let isSearching = false;
    let searchResult = {};

    $: fields = searchResult?.schema||[];
    $: records = searchResult?.records||[];



     function search() {
        excuteErrorMsg=""

        if(!collection?.rawSql || !collection?.did){
            excuteErrorMsg = "Datasource or sql is invalid";
            return;
        }else{
            excuteErrorMsg = "";
        }

        isSearching = true;

        let request  = ApiClient.send("/api/sql-collections/excute-sql",{
            'method': 'POST',
            'body':   {rawSql:collection?.rawSql,did:collection?.did},
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
        search()?.then(
            ()=>{
                if(excuteErrorMsg){
                   return
                }

                let c = collection?.clone();
                c.schema = searchResult.schema;

                sqlConfirmPanel?.show(c,hasChanges)
            }
        );
    }



    const TAB_FIELDS = "fields";
    const dispatch = createEventDispatcher();

    let sqlCollectionPanel;




    let isSaving = false;
    let confirmClose = false; // prevent close recursion
    let activeTab = TAB_FIELDS;
    let initialFormHash = calculateFormHash(collection);



    $: isSystemUpdate = !collection.isNew && collection.system;

    $: hasChanges = initialFormHash != calculateFormHash(collection);

    $: canAsCollection = collection.isNew || hasChanges;

    export function changeTab(newTab) {
        activeTab = newTab;
    }

    export function show(model) {

        load(model);

        confirmClose = true;

        changeTab(TAB_FIELDS);

        if (collection.isNew) {
            sqlCollectionPanel?.show()
        }else{
            sqlCollectionPanel?.show()
            setTimeout(()=>{ // 增加延时，否则两个Panel的z-index相同，导致编辑数据后 cancel弹出框在Pannel后面
                sqlConfirmPanel?.show(collection)
            },1)

        }

    }

    export function hide() {
        return sqlCollectionPanel?.hide();
    }


    async function load(model) {
        setErrors({}); // reset errors
        if (typeof model !== "undefined") {
            original = model;
            collection = model?.clone();
        } else {
            original = null;
            collection = new CollectionExp();
        }

        // normalize
        collection.schema = collection.schema || [];
        collection.originalName = collection.name || "";

        await tick();

        initialFormHash = calculateFormHash(collection);

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
        class="overlay-panel-xxl colored-header compact-header collection-panel"
        beforeHide={() => {
        if (hasChanges && confirmClose) {

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
                            name="datasource"
                            let:uniqueId
                    >
                        <label for={uniqueId}>Datasource</label>
                        <!-- svelte-ignore a11y-autofocus -->
                        <DataSourceSelect bind:keyOfSelected = {collection.did} ></DataSourceSelect>

                    </Field>
                </div>
            </div>

            <div class="grid m-b-base">
                <div class="col-lg-10">

                    <Field class="form-field required m-b-0 {isSystemUpdate ? 'disabled' : ''}"
                           name="rawSql"
                           let:uniqueId>
                        <label for={uniqueId}>Sql Editor</label>
                        <textarea
                                id={uniqueId}
                                class="txt-mono"
                                spellcheck="false"
                                rows="8"
                                required
                                bind:value={collection.rawSql}

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
                disabled={!canAsCollection || isSearching}
                on:click={() => asCollection()}
        >
            <span class="txt">{collection.isNew ? "As Collection" : "Save changes"}</span>
        </button>
    </svelte:fragment>
</OverlayPanel>

<SqlConfirmPanel bind:this={sqlConfirmPanel} on:save={()=>{
    confirmClose = false;
    hide();
}} on:delete={()=>{
    confirmClose = false;
    hide();
}}/>
<style>

</style>
