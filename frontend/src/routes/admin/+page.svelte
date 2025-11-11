<script lang="ts">
    import {
        Button,
        Alert,
        FormGroup,
        Input,
        ButtonGroup,
    } from "@sveltestrap/sveltestrap";
    import * as m from "$lib/paraglide/messages";
    import type { PageProps } from "./$types";
    import { enhance } from "$app/forms";
    import { getLocale } from "$lib/paraglide/runtime.js";
    let { data, form } = $props();
    let isLogin = $state(true);
    let submitInProgress = $state(false);
    let errors = { email: "" };
</script>

{#if data.loggedIn}
    <main class="container">
        <p>admin logged in</p>
    </main>
{:else}
    <main class="container-md login">
        {#if form?.error}
            <Alert dismissible color="danger">
                <h4 class="alert-heading text-capitalize">{m.success()}</h4>
                {form.error}
            </Alert>
        {:else if form?.success}
            <Alert dismissible color="success">
                <h4 class="alert-heading text-capitalize">{m.error()}</h4>
                {form.success}
            </Alert>
        {/if}

        <h2 class="text-center">{m.app_navbar_login()}</h2>
        <form
            action="?/login"
            method="POST"
            use:enhance={() => {
                submitInProgress = true;

                return async ({ update }) => {
                    await update();
                    submitInProgress = false;
                };
            }}
            data-sveltekit-reload
        >
            <FormGroup floating label={m.username()}>
                <Input
                    autocomplete="username"
                    disabled={submitInProgress}
                    name="username"
                    required
                />
            </FormGroup>
            <FormGroup floating label={m.form_login_password()}>
                <Input
                    autocomplete="current-password"
                    disabled={submitInProgress}
                    type="password"
                    name="password"
                    required
                />
            </FormGroup>
            <Button
                block
                disabled={submitInProgress}
                data-sveltekit-reload
                type="submit">{m.app_navbar_login()}</Button
            >
        </form>
    </main>
{/if}
