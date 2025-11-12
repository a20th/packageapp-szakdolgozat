<script lang="ts">
import * as m from "$lib/paraglide/messages";
import {Accordion, AccordionItem} from "@sveltestrap/sveltestrap";

</script>

<main class="container content p-3">
    <h2 class="text-center mb-3">{m.support()}</h2>
    <p>{m.faqlong()}</p>
    <p>{m.phone_number()}: +36101234567</p>
    <p>Email: support@package.hu</p>
    <h2 class="text-center mb-3">{m.app_navbar_faq()} - {m.faqlong()}</h2>
    <Accordion stayOpen>
        <AccordionItem header="{m.question()} #1">
            <p>{m.answer()} #1</p>
        </AccordionItem>
        <AccordionItem header="{m.question()} #2">
            <p>{m.answer()} #2</p>
        </AccordionItem>
        <AccordionItem header="{m.question()} #3">
            <p>{m.answer()} #3</p>
        </AccordionItem>
        <AccordionItem header="{m.question()} #4">
            <p>{m.answer()} #4</p>
        </AccordionItem>
    </Accordion>
</main>
