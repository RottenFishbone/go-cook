<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
	import { apiRoot, State } from '../common';
	import RecipeItem from './recipe_item.svelte';

	export let pageSize: number = 10;
	export let page: number = 0;
	
	const dispatch = createEventDispatcher();

	$: recipes = []
	let searchText = "";
	$: searchText, searchChange();
	
  // Loads the recipe names from the API server
	async function fetchRecipeList() {
		let query = searchText == '' ? '' : `&q=${searchText.replaceAll(' ', ',')}`;
		const resp = await fetch(`${apiRoot}/recipes/names?count=${pageSize}&page=${page}${query}`);
    if (resp.ok) {
      return resp.json();
    } else {
			throw new Error(`Failed to fetch recipe list: ${resp.status} ${resp.statusText}`);
    }
  }

	let failedLoad = false;
	let fetched = false;
  let mounted = false;

  // Load the recipes immediately
  onMount(async () => {
    // This timer will hide content briefly to avoid flashing the user with text
    setTimeout(()=>{
      mounted = true;
    }, 200);

    try {
			recipes = await fetchRecipeList();
			fetched = true;
    } catch (err) {
      failedLoad = true;
      throw err
    }
	});

	// Handles on click event for `New Recipe` button
	function clickNewRecipe() {
		handleItmMsg({ detail: {
			tag: 'new',
			msg: '',
		}});
	}

	let debouncing: boolean = false;
	let searchPending: boolean = false;
	// A debounced recipe search
	async function searchChange() {
		// Prevent rapid calling while the text is changing quickly
		if (debouncing && searchText != '') {
			searchPending = true;
			return;
		}

		// Prevent rapid calls and reset the page (for new results)
		debouncing = true;
		page = 0;
		try {
			recipes = await fetchRecipeList()
		} catch (err) {
			failedLoad = true;
			throw err;
		}

		// Set a delay to unset `debounce` and allow searching again
		setTimeout(async ()=>{
			debouncing = false;

			// If a search was attempted during debounce, we'll call it now
			if (searchPending) {
				searchPending = false;
				await searchChange()
			}
		}, 250);
	}

  // Handles an event thrown by a RecipeItem
  function handleItmMsg(event: { detail: { tag: string, msg: string } }) {
    let msg = event.detail.msg;
    switch (event.detail.tag) {
      case 'delete':
				fetch(`${apiRoot}/recipes/?name=${msg}`, {
          method: 'DELETE',
        }).then(resp => {
          if (resp.ok){
            let id = recipes.findIndex((v=>v==msg));
            recipes.splice(id, 1);
            recipes = recipes;
          } else {
            // TODO convert to a toast
            console.log(`Recipe delete rejected: ${resp.status} ${resp.statusText}`);
          }
        }).catch((err) => {
          // TODO convert to a toast
					alert('Failed to connect to API server: ' + err);
        });
        break;
		 	case 'edit':
        dispatch('msg', {
          tag: State.RecipeEdit,
          msg: event.detail.msg,
        });
        break;
      case 'view':
        // Bubble the event up to a component that can handle it
        dispatch('msg', {
          tag: State.RecipeView,
          msg: event.detail.msg,
        });
				break;
 			case 'new':
				dispatch('msg', {
					tag: State.RecipeEdit,
					msg: '',
				});
				break;
      default:
				console.log('Unknown message recieved from recipe_item component.');
    }
  }
</script>

<!-- Main Recipe List -->
{#if fetched}
<div class="my-2 rounded-box flex-col">
  <!-- Search bar -->
  <div class="flex justify-center my-2">
		<input 
			 bind:value={searchText}
       type="text" 
       placeholder="Search Recipes..." 
       class="input input-bordered input-sm w-full max-w-md"/>
  </div>
  <!-- Recipe List -->
  <div class="flex justify-center">
		<ul class="lower-z rounded-box max-w-md w-full p-2">
			<li class="my-2 flex">
				<button class="btn btn-outline btn-primary normal-case flex-1" on:click={clickNewRecipe}>
					New Recipe
				</button>
			</li>
			{#each recipes as recipe (recipe)}
				<li class="my-2"><RecipeItem recipeName={recipe} on:msg={handleItmMsg}/></li>
			{/each}
    </ul>
  </div>
</div>

<!-- Loading spinner -->
{:else if !failedLoad}
	<div class={`flex flex-col justify-center mx-auto max-w-md transition-opacity duration-750 
						 ${mounted ? '' : 'min-h-screen opacity-0'}`}>
    <div class="text-xl mx-auto my-5">Fetching recipes...</div>
    <div class="btn btn-circle btn-xl btn-disabled mx-auto loading btn-primary"></div>
  </div>
{:else}
  <div class="flex justify-center mx-auto max-w-md">
    <div class="my-5 bg-warning text-warning-content text-lg p-2">
      Failed to load Recipe List from server.<br/>
      Try Refreshing?
    </div>
  </div>
{/if}
