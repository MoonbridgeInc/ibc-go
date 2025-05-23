---
title: Governance Proposals
sidebar_label: Light Client Recovery
sidebar_position: 6
slug: /ibc/proposals
---

# Light Client Recovery

In uncommon situations, a highly valued client may become frozen or expire due to uncontrollable
circumstances. A highly valued client might have hundreds of channels being actively used.
Some of those channels might have a significant amount of locked tokens used for ICS 20.

## Frozen Light Clients

If the one third of the validator set of the chain the client represents decides to collude,
they can sign off on two valid but conflicting headers each signed by the other one third
of the honest validator set. The light client can now be updated with two valid, but conflicting
headers at the same height. The light client cannot know which header is trustworthy and therefore
evidence of such misbehaviour is likely to be submitted resulting in a frozen light client.

Frozen light clients cannot be updated under any circumstance except via a governance proposal.
Since a quorum of validators can sign arbitrary state roots which may not be valid executions
of the state machine, a governance proposal has been added to ease the complexity of unfreezing
or updating clients which have become "stuck". Without this mechanism, validator sets would need
to construct a state root to unfreeze the client. Unfreezing clients, re-enables all of the channels
built upon that client. This may result in recovery of otherwise lost funds.

## Expired Light Clients

Tendermint light clients may become expired if the trusting period has passed since their
last update. This may occur if relayers stop submitting headers to update the clients.

An unplanned upgrade by the counterparty chain may also result in expired clients. If the counterparty
chain undergoes an unplanned upgrade, there may be no commitment to that upgrade signed by the validator
set before the chain ID changes. In this situation, the validator set of the last valid update for the
light client is never expected to produce another valid header since the chain ID has changed, which will
ultimately lead the on-chain light client to become expired.

# How to recover an expired client with a governance proposal

> **Who is this information for?**
> Although technically anyone can submit the governance proposal to recover an expired client, often it will be **relayer operators** (at least coordinating the submission).

In the case that a highly valued light client is frozen, expired, or rendered non-updateable, a
governance proposal may be submitted to update this client, known as the subject client. The
proposal includes the client identifier for the subject and the client identifier for a substitute
client. Light client implementations may implement custom updating logic, but in most cases,
the subject will be updated to the latest consensus state of the substitute client, if the proposal passes.
The substitute client is used as a "stand in" while the subject is on trial. It is best practice to create
a substitute client *after* the subject has become frozen to avoid the substitute from also becoming frozen.
An active substitute client allows headers to be submitted during the voting period to prevent accidental expiry
once the proposal passes.

See also the relevant documentation: [ADR-026, IBC client recovery mechanisms](/architecture/adr-026-ibc-client-recovery-mechanisms)

## Preconditions

- There exists an active client (with a known client identifier) for the same counterparty chain as the expired client.
- The governance deposit.

## Steps

### Step 1

Check if the client is attached to the expected `chain_id`. For example, for an expired Tendermint client representing the Akash chain the client state looks like this on querying the client state:

```text
{
  client_id: 07-tendermint-146
  client_state:
  '@type': /ibc.lightclients.tendermint.v1.ClientState
  allow_update_after_expiry: true
  allow_update_after_misbehaviour: true
  chain_id: akashnet-2
}
```

The client is attached to the expected Akash `chain_id`. Note that although the parameters (`allow_update_after_expiry` and `allow_update_after_misbehaviour`) exist to signal intent, these parameters have been deprecated and will not enforce any checks on the revival of client. See ADR-026 for more context on this deprecation.

### Step 2

Anyone can submit the governance proposal to recover the client by executing the following via CLI.
If the chain is on an ibc-go version older than v8, please see the [relevant documentation](https://ibc.cosmos.network/v7/ibc/proposals).

- From ibc-go v8 onwards

  ```shell
  <binary> tx gov submit-proposal [path-to-proposal-json]
  ```

  where `proposal.json` contains:

  ```json
  {
    "messages": [
      {
        "@type": "/ibc.core.client.v1.MsgRecoverClient",
        "subject_client_id": "<expired-client-id>",
        "substitute_client_id": "<active-client-id>",
        "signer": "<gov-address>"
      }
    ],
    "metadata": "<metadata>",
    "deposit": "10stake"
    "title": "My proposal",
    "summary": "A short summary of my proposal",
    "expedited": false
  }
  ```

The `<expired-client-id>` identifier is the proposed client to be updated. This client must be either frozen or expired.

The `<active-client-id>` represents a substitute client. It carries all the state for the client which may be updated. It must have identical client and chain parameters to the client which may be updated (except for latest height, frozen height, and chain ID). It should be continually updated during the voting period.

After this, all that remains is deciding who funds the governance deposit and ensuring the governance proposal passes. If it does, the client on trial will be updated to the latest state of the substitute.

## Important considerations

Please note that if the counterparty client is also expired, that client will also need to update. This process updates only one client.
