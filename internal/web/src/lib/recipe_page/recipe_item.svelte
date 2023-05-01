<script lang="ts">
  import {createEventDispatcher} from 'svelte';
  
  // Recipe title
  export let title: string;

  const DELETE_DEFAULT = "Delete";
  const DELETE_CONFIRM = "You Sure?";

  const dispatch = createEventDispatcher();

  let focused = false;
  let secondPress = false;
  let deleteText = DELETE_DEFAULT;
  
  // component's DOM binding
  let group: HTMLDivElement;

  // Handle focusin on button group
  function focusin(_: FocusEvent) {
    focused = true;
  }

  // Handle focusout on button group
  function focusout(e: FocusEvent) {
    // Only focus out if the new focus is outside the button group
    let target = e.relatedTarget;
    if (!group.contains(target as Node)){
      focused = false;
      secondPress = false;
    } 
    deleteText = DELETE_DEFAULT;
  }

  function recipeClick() {
    if (!focused) { return; }
    if (!secondPress) { secondPress = true; return; }

    //TODO fire view recipe event
  }

  function deleteClick() {
    if (deleteText==DELETE_DEFAULT) { 
      deleteText=DELETE_CONFIRM; 
      return; 
    }

    dispatch('msg', {
      tag: 'delete',
      msg: title,
    })
  }

  function editClick() {
    //TODO fire edit recipe event
  }


</script>

<div bind:this={group} class="btn-group flex flex-row" on:focusin={focusin} on:focusout={focusout}>
  <!-- On focused state, allow editng and recipe clicking -->
  <button class="flex-1 btn normal-case" on:click={recipeClick}>
    {title}
  </button>
  <div class="dropdown dropdown-left {focused ? '' : 'hidden'}">
    <label tabindex="-1" class="btn m-1 text-xl">+</label>
    <ul tabindex="-1" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-40">
      <li class="my-1"><a class="justify-center">Edit Recipe</a></li>
      <li class="my-1">
        <a class={`justify-center
          bg-warning text-warning-content 
          hover:bg-primary hover:text-primary-content`}
           on:click={deleteClick}>
          {deleteText}
        </a>
      </li>
    </ul>
  </div>
</div>

<style>
</style>
