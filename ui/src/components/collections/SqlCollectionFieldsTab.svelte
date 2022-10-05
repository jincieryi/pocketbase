<script>


    export let collection = {};

    $: if (typeof collection?.schema === "undefined") {
        collection = collection || {};
        collection.schema = [];
    }


</script>


<div class="table-wrapper">
    <table class="table">
        <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>As Id</th>

        </tr>
        </thead>
        <tbody>
        {#each collection.schema as field (field.name)}
          <tr><td>{field.name}</td><td>{field.type}</td><td></td></tr>
        {/each}

        </tbody>
    </table>
</div>

<div class="clearfix m-t-xs" />

