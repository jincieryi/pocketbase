<script>
    import { replace, querystring } from "svelte-spa-router";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle } from "@/stores/app";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import Searchbar from "@/components/base/Searchbar.svelte";
    import RefreshButton from "@/components/base/RefreshButton.svelte";
    import SortHeader from "@/components/base/SortHeader.svelte";
    import IdLabel from "@/components/base/IdLabel.svelte";
    import FormattedDate from "@/components/base/FormattedDate.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import DatasourceUpsertPanel from "@/components/datasources/DatasourceUpsertPanel.svelte";

    $pageTitle = "Datasources";

    const queryParams = new URLSearchParams($querystring);

    let datasourceUpsertPanel;
    let datasources = [];
    let isLoading = false;
    let filter = queryParams.get("filter") || "";
    let sort = queryParams.get("sort") || "-created";

    $: if (sort !== -1 && filter !== -1) {
        // keep listing params in sync
        const query = new URLSearchParams({ filter, sort }).toString();
        replace("/settings/datasources?" + query);

        loadDatasources();
    }

    export function loadDatasources() {
        isLoading = true;

        datasources = []; // reset

        return ApiClient.records
            .getList("datasources",1,100, {
                sort: sort || "-created",
                filter: filter,
            })
            .then((result) => {
                datasources = result.items;
                isLoading = false;
            })
            .catch((err) => {
                if (!err?.isAbort) {
                    isLoading = false;
                    console.warn(err);
                    clearList();
                    ApiClient.errorResponseHandler(err, false);
                }
            });
    }

    function clearList() {
        datasources = [];
    }
</script>

<SettingsSidebar />

<PageWrapper>
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>

        <RefreshButton on:refresh={() => loadDatasources()} />

        <div class="flex-fill" />

        <button type="button" class="btn btn-expanded" on:click={() => datasourceUpsertPanel?.show()}>
            <i class="ri-add-line" />
            <span class="txt">New datasource</span>
        </button>
    </header>

    <Searchbar
        value={filter}
        placeholder={"Search filter, eg. name='myprd'"}
        extraAutocompleteKeys={["name"]}
        on:submit={(e) => (filter = e.detail)}
    />

    <div class="table-wrapper">
        <table class="table" class:table-loading={isLoading}>
            <thead>
                <tr>
                    <th class="min-width" />

                    <SortHeader class="col-type-text" name="id" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("primary")} />
                            <span class="txt">id</span>
                        </div>
                    </SortHeader>

                    <SortHeader class="col-type-text col-field-text" name="name" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("text")} />
                            <span class="txt">name</span>
                        </div>
                    </SortHeader>
                    <SortHeader class="col-type-text col-field-text" name="dsn" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("text")} />
                            <span class="txt">dsn</span>
                        </div>
                    </SortHeader>
                    <SortHeader class="col-type-text col-field-text" name="type" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("text")} />
                            <span class="txt">type</span>
                        </div>
                    </SortHeader>

                    <SortHeader class="col-type-date col-field-created" name="created" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("date")} />
                            <span class="txt">created</span>
                        </div>
                    </SortHeader>

                    <SortHeader class="col-type-date col-field-updated" name="updated" bind:sort>
                        <div class="col-header-content">
                            <i class={CommonHelper.getFieldTypeIcon("date")} />
                            <span class="txt">updated</span>
                        </div>
                    </SortHeader>

                    <th class="col-type-action min-width" />
                </tr>
            </thead>
            <tbody>
                {#each datasources as datasource (datasource.id)}
                    <tr
                        tabindex="0"
                        class="row-handle"
                        on:click={() => datasourceUpsertPanel?.show(datasource)}
                        on:keydown={(e) => {
                            if (e.code === "Enter" || e.code === "Space") {
                                e.preventDefault();
                                datasourceUpsertPanel?.show(datasource);
                            }
                        }}
                    >
                    <td class="min-width"></td>
                        <td class="col-type-text col-field-id">
                            <IdLabel id={datasource.id} />
                        
                        </td>

                        <td class="col-type-text col-field-name">
                            <span class="txt txt-ellipsis" title={datasource.name}>
                                {datasource.name}
                            </span>
                        </td>

                        <td class="col-type-text col-field-dsn">
                            <span class="txt txt-ellipsis" title={datasource.dsn}>
                                {datasource.dsn}
                            </span>
                        </td>

                        <td class="col-type-text col-field-type">
                            <span class="txt txt-ellipsis" title={datasource.type}>
                                {datasource.type}
                            </span>
                        </td>

                        <td class="col-type-date col-field-created">
                            <FormattedDate date={datasource.created} />
                        </td>
                        <td class="col-type-date col-field-updated">
                            <FormattedDate date={datasource.updated} />
                        </td>
                        <td class="col-type-action min-width">
                            <i class="ri-arrow-right-line" />
                        </td>
                    </tr>
                {:else}
                    {#if isLoading}
                        <tr>
                            <td colspan="99" class="p-xs">{isLoading}
                                <span class="skeleton-loader" />
                            </td>
                        </tr>
                    {:else}
                        <tr>
                            <td colspan="99" class="txt-center txt-hint p-xs">
                                <h6>No datasources found.</h6>
                                {#if filter?.length}
                                    <button
                                        type="button"
                                        class="btn btn-hint btn-expanded m-t-sm"
                                        on:click={() => (filter = "")}
                                    >
                                        <span class="txt">Clear filters</span>
                                    </button>
                                {/if}
                            </td>
                        </tr>
                    {/if}
                {/each}
            </tbody>
        </table>
    </div>

    {#if datasources.length}
        <small class="block txt-hint txt-right m-t-sm">Showing {datasources.length} of {datasources.length}</small>
    {/if}
</PageWrapper>

<DatasourceUpsertPanel bind:this={datasourceUpsertPanel} on:save={() => loadDatasources()} on:delete={() => loadDatasources()} />
