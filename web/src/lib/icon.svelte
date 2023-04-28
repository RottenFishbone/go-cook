<script lang="ts">
  export let name: string;
  export let width: string;
  export let height: string;

  let icoPath = "src/lib/icons/"
</script>

<svg class={$$props.class} {width} {height}>
  <use href={`${icoPath}${name}.svg#${name}`} />
</svg>
