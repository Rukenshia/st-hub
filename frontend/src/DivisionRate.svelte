<script>
    export let name;
    export let battles;
    export let division;
    export let rate;
</script>

<style>
    .good {
        color: #5ab15a;
    }

    .bad {
        color: #da3c3c;
    }

    .rate {
        font-weight: 500;
    }

    td {
        padding-top: 8px;
        padding-right: 16px;
    }

    td + td {
        text-align: right;
    }
</style>

<tr>
    <td>{name}</td>
    <td><span class="rate" class:good={rate < 50} class:bad={rate >= 50}>{rate}%</span></td>
    <td>{division} / {battles}</td>
</tr>