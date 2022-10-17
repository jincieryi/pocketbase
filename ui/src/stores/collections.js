import {writable} from "svelte/store";
import ApiClient from "@/utils/ApiClient";
import CommonHelper from "@/utils/CommonHelper";
import {Collection} from "pocketbase";

export const collections = writable([]);
export const activeCollection = writable({});
export const isCollectionsLoading = writable(false);


// add or update collection
export function addCollection(collection) {
    activeCollection.update((current) => {
        return CommonHelper.isEmpty(current?.id) || current.id === collection.id ? collection : current;
    });

    collections.update((list) => {
        CommonHelper.pushOrReplaceByKey(list, collection, "id");
        return list;
    });
}

export function removeCollection(collection) {
    collections.update((list) => {
        CommonHelper.removeByKey(list, "id", collection.id);

        activeCollection.update((current) => {
            if (current.id === collection.id) {
                // fallback to the first non-profile collection item
                return list.find((c) => c.name != import.meta.env.PB_PROFILE_COLLECTION) || {}
            }
            return current;
        });

        return list;
    });

}

/**
 * CollectionExp  Collection的扩展
 */
export class CollectionExp extends Collection {
    constructor() {
        super();
    }

    /**
     * Exports all model properties as a new plain object.
     */
    export() {
        return Object.assign({}, this);
    }

    load(data) {

        super.load(data);
        this.setExp(data)
    }

    setExp(data={}){
        this.rawSql = typeof data.rawSql === 'string' ? data.rawSql : null;
        this.did = typeof data.did === 'string' ? data.did : null;
        this.cid = typeof data.cid === 'string' ? data.cid : null;
    }

    get isSqlType(){
        return typeof this.rawSql === 'string' && this.rawSql != ""
    }

    clone() {
        let cloneObj = new CollectionExp();
        Object.assign(cloneObj, JSON.parse(JSON.stringify(this)));
        return cloneObj;
    }

}

// load all collections (excluding the user profile)
export async function loadCollections(activeId = null) {
    isCollectionsLoading.set(true);

    activeCollection.set({});
    collections.set([]);


    return ApiClient.collections.getFullList(200, {
        "sort": "+created",

    }).then(async (collections) => {
        // 拼接collectionsExtend 过滤条件
        let cids = [];
        let items = [] // item type is CollectionExp

        collections.forEach((item) => {
            let collection = new CollectionExp();
            collection.load(item)
            items.push(collection)

            cids.push("cid=\"" + item.id + "\"");
        });

        let recorders = await ApiClient.records.getFullList(
            "collectionsExtend",
            200,
            {"filter": CommonHelper.joinNonEmpty(cids, " || ")}
        );

        items.forEach((item) => {
            let exp = CommonHelper.findByKey(recorders, "cid", item.id)
            if (!CommonHelper.isEmpty(exp) ){
                item.setExp(exp)
            }

        });


        return items;
    })
        .then((items) => {
            collections.set(items);

            const item = activeId && CommonHelper.findByKey(items, "id", activeId);
            if (item) {
                activeCollection.set(item);
            } else if (items.length) {
                // fallback to the first non-profile collection item
                const nonProfile = items.find((c) =>
                    (c.name != import.meta.env.PB_PROFILE_COLLECTION
                        && c.name != import.meta.env.PB_DATASOURCE_COLLECTION))
                if (nonProfile) {
                    activeCollection.set(nonProfile);
                }
            }

        })
        .catch((err) => {
            ApiClient.errorResponseHandler(err);
            console.error(err)
        })
        .finally(() => {
            isCollectionsLoading.set(false);
        });
}
