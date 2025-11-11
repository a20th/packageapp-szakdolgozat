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
    let id = $state();
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
    <h2 class="text-center mb-4">Admin users</h2>

    <!-- svelte-ignore component_name_lowercase -->
    <form
        action="?/create"
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
        <h3 class="text-center mb-4">Create new user</h3>
        <div class="row">
            <div class="col">
                <FormGroup floating>
                    <Input disabled={submitInProgress} name="username" required
                    ></Input>
                    <div slot="label">Username</div>
                </FormGroup>
            </div>
            <div class="col">
                <FormGroup floating>
                    <Input
                        disabled={submitInProgress}
                        name="Password"
                        type="password"
                        required
                    />
                    <div slot="label">Password</div>
                </FormGroup>
            </div>
            
        </div>
        <Button type="submit" class="w-100" disabled={submitInProgress}
                >Create</Button
            >
    </form>
    <hr />
    <h3 class="text-center mb-4">Manage users</h3>
    
        {#each data.admins as user}
                <form
                    action="?/delete"
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
                    <hr class="w-75 m-auto my-4"/>
                    <div class="row">
                        <div class="col">
                            <h4 class="text-center">{user}</h4>
                        </div>
                        <input type="hidden" value={user} name="username" />
                        <div class="col">
                            <Button
                                type="submit"
                                class=""
                                color="danger"
                                disabled={submitInProgress || user == "admin"}>Delete</Button
                            >
                        </div>
                    </div>
                    <hr class="w-75 m-auto my-4"/>
                </form>

        {/each}

</main>
