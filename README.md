# Apollo Backend — Self-Hosting Guide

This guide covers running the Apollo push notification backend yourself, paired with the [Apollo-ImprovedCustomAPI](https://github.com/JeffreyCA/Apollo-ImprovedCustomAPI) tweak which lets you point the Apollo app at your own server.

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running
- A Reddit account and a Reddit API key (created at https://www.reddit.com/prefs/apps)
  - **NOTE**: Due to the new [Responsible Builder Policy](https://www.reddit.com/r/redditdev/comments/1oug31u/introducing_the_responsible_builder_policy_new/), if you have not already created an API key, you cannot do so anymore. See https://github.com/JeffreyCA/Apollo-ImprovedCustomApi?tab=readme-ov-file#dont-have-an-api-key for workarounds to get an API key.
- An [Apple Developer Program](https://developer.apple.com/programs/) membership ($99/year) — required for push notifications
- The Apollo-ImprovedCustomAPI tweak installed on your device
  - Use [my fork](https://github.com/DeltAndy123/Apollo-ImprovedCustomApi) as it has the option to configure push notification server URL

---

## Step 1 — Clone and Configure

```bash
git clone https://github.com/DeltAndy123/apollo-backend
cd apollo-backend
cp .env.example .env
```

Open `.env` and fill in all the values. See the sections below for where to find each one.

---

## Step 2 — Apple APNs Key (required for push notifications)

Push notifications are delivered via Apple's Push Notification service (APNs). You need an Apple Developer account and an APNs key tied to a bundle ID you control.

> [!TIP]
> Don't have an Apple Developer account yet? You can generate fake credentials to get the server running for testing. The server will start and all features except push notifications will work. See [Fake APNs credentials for testing](#fake-apns-credentials-for-testing) below.

### 2a — Register a Bundle ID

1. Go to [developer.apple.com](https://developer.apple.com) → Certificates, Identifiers & Profiles → Identifiers
2. Click **+** → **App IDs** → **App**
3. Enter a description and a Bundle ID (e.g. `com.yourname.Apollo`)
4. Under Capabilities, enable **Push Notifications**
5. Click Continue → Register

### 2b — Create an APNs Key

1. Go to Certificates, Identifiers & Profiles → **Keys**
2. Click **+** → check **Apple Push Notifications service (APNs)**
3. Download the `.p8` file — **you can only download it once**
4. Note the **Key ID** shown on the confirmation page (10 characters)
5. Your **Team ID** is shown in the top-right of the developer portal, or under Membership

### 2c — Configure `.env`

```env
APPLE_BUNDLE_ID=com.yourname.Apollo
APPLE_KEY_ID=XXXXXXXXXX
APPLE_TEAM_ID=XXXXXXXXXX
APPLE_KEY_PATH=/etc/secrets/apple.p8
```

### 2d — Place the `.p8` file

```bash
mkdir -p secrets
cp ~/Downloads/AuthKey_XXXXXXXXXX.p8 secrets/apple.p8
```

Docker Compose mounts `./secrets` into the container at `/etc/secrets`.

---

## Step 3 — Resign and Install Apollo

Because your APNs key is tied to your bundle ID (not `com.christianselig.Apollo`), you need to resign the Apollo IPA with your bundle ID before installing.

Tools that can do this:
- **[Feather](https://feather.khcrysalis.dev/) or [Impactor](https://impactor.khcrysalis.dev/)**: Change the "Identifier" when signing Apollo
- **[Sideloadly](https://sideloadly.io/)**: Enable "Change app bundle ID" in Advanced options

The bundle ID you set here must match `APPLE_BUNDLE_ID` in `.env`.

> [!IMPORTANT]
> Free Apple Developer accounts cannot be used for push notifications. Anyone who sideloads the IPA must have a paid Apple Developer Program membership, or must have access to a signing service such as [ArcticSign](https://arcticsign.app/) to use push notifications.

---

## Step 4 — Start the Server

```bash
docker compose up --build -d
```

Check that the API is running:

```bash
curl http://localhost:4000/v1/health
# {"status":"available"}

curl http://localhost:4000/v1/bundle_id
# {"bundle_id":"com.yourname.Apollo"}
```

To watch logs:

```bash
docker compose logs -f api
```

---

## Step 5 — Configure Apollo

In Apollo → Settings → Custom API → Scroll to Push Notifications:

- **Push Server URL**: `http://your-server-ip:4000` (or your domain if publicly hosted)
- **Reddit Client ID**: the same client ID from Step 2

Then go to Settings → Notifications. Apollo will register your device and link your Reddit account.

---

## Fake APNs credentials for testing

> [!WARNING]
> With fake credentials the server will start and accept connections, but **push notifications will not be delivered**. Watchers, account linking, and all other API features will still work. Use this only to test the setup before obtaining real Apple Developer credentials.

The server requires a valid EC private key file on startup — it won't accept an empty file. Generate one with OpenSSL (installed by default on macOS and most Linux distros):

```bash
mkdir -p secrets
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out secrets/apple.p8
```

Then set placeholder values for the Key ID and Team ID in `.env` (they must be exactly 10 characters):

```env
APPLE_BUNDLE_ID=com.yourname.Apollo
APPLE_KEY_ID=FAKEKEYID0
APPLE_TEAM_ID=FAKETEAMID
APPLE_KEY_PATH=/etc/secrets/apple.p8
```

The server will start successfully. When a notification would be sent, Apple will reject it with `InvalidProviderToken` — this will appear in the logs but won't crash the server.

---

## Useful Commands

```bash
# Start everything
docker compose up -d

# Rebuild after code changes
docker compose up --build -d api

# View logs
docker compose logs -f api
docker compose logs -f worker-notifications

# Connect to the database
docker compose exec postgres psql -U apollo apollo

# Stop everything
docker compose down

# Stop and delete all data
docker compose down -v
```

---

## How It Works

| Service | Role |
|---------|------|
| `api` | HTTP API on port 4000 — handles device registration, account linking, receipt validation, watchers |
| `scheduler` | Enqueues notification-check jobs every 5 seconds |
| `worker-notifications` | Polls Reddit inboxes and sends APNs pushes |
| `worker-users` | User activity watchers |
| `worker-subreddits` | Subreddit keyword/flair/author watchers |
| `worker-trending` | Trending post watchers |
| `worker-stuck` | Recovers jobs stuck in the queue |
| `postgres` | Persistent storage for accounts, devices, watchers |
| `redis-queues` | Job queue (rmq) |
| `redis-locks` | Distributed locks to prevent duplicate notifications |

Receipt validation is bypassed on a self-hosted server — the `/v1/receipt` endpoint always returns a lifetime subscription for all products, since sideloaded apps cannot use App Store in-app purchases.

---

---

# Original README

> *The following is the original README from when this repository was open-sourced.*

## Apollo Server Backend

This repository holds Apollo's code for its server backend, which checks for user notifications and allows users to create subreddit watchers. It is archived as it will no longer be used/updated after June 30th, 2023 per Reddit's decisions in regards to the Reddit API. https://www.reddit.com/r/apolloapp/comments/144f6xm/apollo_will_close_down_on_june_30th_reddits/

The goal of making the code for this repo available is to show that despite statements otherwise by Reddit administrators, Apollo does not scrape anything and users purely authenticated Reddit API requests, and does a great deal of work to ensure the Reddit API rate limits are respected.
