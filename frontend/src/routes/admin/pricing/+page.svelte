<script lang="ts">
    import { enhance } from "$app/forms";
    import * as m from "$lib/paraglide/messages";
    import {
        Button,
        Alert,
        FormGroup,
        Input,
        ButtonGroup,
    } from "@sveltestrap/sveltestrap";
    let { data, form } = $props();
    let submitInProgress = $state(false);
</script>

{#if form?.success}
    <Alert class="container content" dismissible id="error" color="success">
        <h4 class="alert-heading text-capitalize">{m.success()}</h4>
    </Alert>
{:else if form?.success === false}
    <Alert class="container content" dismissible id="error" color="danger">
        <h4 class="alert-heading text-capitalize">{m.error()}</h4>
    </Alert>
{/if}

<main class="container content p-4">
    <h2 class="text-center mb-4">{m.setpricing()}</h2>

    <!-- svelte-ignore component_name_lowercase -->
    <form
        action=""
        method="post"
        class=""
        use:enhance={() => {
            submitInProgress = true;

            return async ({ update }) => {
                submitInProgress = false;
                await update({ reset: false, invalidateAll: true });
            };
        }}
    >
        <hr />
        <div class="row">
            <div class="col">
                <FormGroup floating>
                    <Input
                        disabled={submitInProgress}
                        bind:value={data.pricing.baseprice}
                        name="baseprice"
                        required
                    ></Input>
                    <div slot="label">{m.baseprice()}</div>
                </FormGroup>
            </div>
            <div class="col">
                <FormGroup floating>
                    <Input
                        disabled={submitInProgress}
                        name="kmprice"
                        type="number"
                        bind:value={data.pricing.kmprice}
                        required
                    />
                    <div slot="label">{m.kmprice()}</div>
                </FormGroup>
            </div>
        </div>
        <Button type="submit" class="w-100" disabled={submitInProgress}
            >{m.update()}</Button
        >
    </form>
</main>
