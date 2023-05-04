<script lang="ts">
  import { stripRecipeName } from '../common';
  import {createEventDispatcher} from 'svelte';
  
  // Recipe title
  export let recipeName: string;
  let title = stripRecipeName(recipeName);

  const dispatch = createEventDispatcher();

  const DELETE_DEFAULT = "Delete";    // Default content of delete button
  const DELETE_CONFIRM = "You Sure?"; // Content to display on first delete press
  let deleteText = DELETE_DEFAULT; 


  let focused = false;      // Tracks if the item is currently focused
  let delPressed = false;   // Tracks if delete has been pressed (for confirms)
  
  // This component's outermost DOM binding
  let group: HTMLDivElement;

  // Handle `focusin` on button group
  function focusin(_: FocusEvent) {
    focused = true;
  }

  // Handle `focusout` on button group
  function focusout(e: FocusEvent) {
    // Only focus out if the new focused element is outside this component
    let target = e.relatedTarget;
    if (!group.contains(target as Node)){
      focused = false;
      delPressed = false;
      deleteText = DELETE_DEFAULT;
    } 
  }

  // Recipe title button handler
  function recipeClick() {
    if (!focused) { return; }
    if (!delPressed) { delPressed = true; return; }
    dispatch('msg', {
      tag: 'view',
      msg: recipeName
    });
  } 

  // Delete button handler
  function deleteClick() {
    if (deleteText==DELETE_DEFAULT) { 
      deleteText=DELETE_CONFIRM; 
      return; 
    }

    dispatch('msg', {
      tag: 'delete',
      msg: recipeName,
    });
  }

  // Edit button handler
  function editClick() {
    dispatch('msg', {
      tag: 'edit',
      msg: recipeName,
    });
  }
</script>

<div bind:this={group} class="flex flex-row" on:focusin={focusin} on:focusout={focusout}>
  <!-- On focused state, allow editng and recipe clicking -->
  <button class="flex-1 btn normal-case" on:click={recipeClick}>
    {title}
  </button>
  
  <!-- Dropdown and expand button -->
  <div class="dropdown dropdown-left {focused ? '' : 'hidden'} flex-none">
    <button tabindex="-1" class="text-xl btn btn-square">+</button>
    
    <!-- Dropdown menu -->
    <ul tabindex="-1" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-40">
      <!-- `Edit` button -->
      <li class="my-1">
        <button class="text-neutral-content bg-neutral justify-center btn-disabled" on:click={editClick}>
          Edit Recipe (Unimpl)
        </button>
      </li>
      <!-- Delete Button -->
      <li class="my-1">
        <button class={`justify-center
          bg-warning text-warning-content 
          hover:bg-primary hover:text-primary-content`}
           on:click={deleteClick}>
          {deleteText}
        </button>
      </li>

    </ul>
  </div>
</div>

<style>
</style>
