# Chicken Service 🐔  
_Even a chicken will understand_

### What is this?  

Chicken Service is a collection of web services built on different languages  
and frameworks. Each language will have its own branch and idea. For example, in **Golang**, we have something like a messenger. Each version in each language represents a development stage of the application. It all starts with the simplest and stupidest implementation and ends with a full-fledged (or almost) good little demonstration application.  

**So**, if you want to learn a new language and understand the basic principles of making apps in this language, **take one chicken to go!** Or, if you want to develop a service in another language or improve the code that I already made, you're also welcome!  

---

**Disclaimer:** I can screw up. That's normal. If you think I wrote actual shit, please explain my mistake, suggest options for improvement, and write about it. This is for the common good of beginners who will watch the code.  

---

## Golang Chicken Messenger  

So here you go. A small messenger.  

**Second disclaimer:** If you don't know something I'm writing here about, check the links at the end of the README. I will leave references to useful topics and basic concepts that I am talking about.  

What do we got here?  

This is a pretty simple microservice application with two services - an auth service and a messenger service._Of course_, for the first version of the application, we could use a monolith, but let's **create something really cool 🤘** After all, we are strong and smart guys.  

### Auth service  
**A responsible service for making new users, logging in users, and validating tokens.**  

- **Register**  
  Register a new user, save the password as a hash. Why hash? Because if you database gets compromised, it won't be fun if hackers steal the actual passwords of your users.  
  Then, return a token. Pretty simple, isn't it?

- **Login**  
  Log in and also return token.  

Now this is our componets. Lets talk about architecture!

---

### Basic architecture  

I've been thinking for a long time about the architecture and came up with the idea of creating something like a **clean architecture**. What does it mean?  

Every logic layer has its own folder, layer, or code base, call it whatever you want. We have **four** layers:  

- **Domain (Core)**  
- **Repository**  
- **Service**  
- **Handler**  

Why it looks like this? To make the code reusable and to separate the logic, of course!  

#### Layers explained:  
- **Domain**: Defines the data (actual database columns), types, and auxiliary structures.
- **Repository**: Defines how data is mapped to the actual database, and what fields does it have (we have 
repositories for PostgreSQL and Redis). Handles the creation of data and stores it in the database 
- **Service**: Defines how data is passed to the repository. Here, we ensure that, for example, there are no duplicate records in the database and prepare the data before transferring it.  
- **Handler**: Handles data from requests and passes it to the service. It checks whether the incoming data is valid.

#### Summary of responsibilities:  
- **Handler**: Reads data from the request, checks types, and "stupidly" returns data or `ApiError` type from service. Check this in code!.  
- **Service**: Validates ALL data. This is one of smartest things in this chain. Its also handling special errors from the database and returns human-like responses.
- **Repository**: Handles all database operations.  
- **Domain**: The data model and core of the application.  

Go to the code and check this out!  

---

### Authorization  


Basic **JWT authorization.** Register a user (having previously checked the input data) and create JWT token. (In this version, I'm using only one token, without the access and refresh concept). So, one token - one session.

But it's not as simple as it might seem in authorization. You can ask: **why do we need Redis here?** Lets talk about token theft

So if token is stolen, we can add endpoint to nullify it. And if we nullify token - it no longer valid. 
And here the question arises: how do we check that a token is valid and not revoked? **Redis!** We save **each token for each session in Redis and validate every request.** If the token is revoked, **bad request, get out of here!**  

**Key format:** `user_id:fingerprint_hash`  
**Fingerprint**: **A string composed of the User-Agent and IP, hashed using murmur3.** **(This is my implementation, you can customize it as you wish).**  

#### Validation process:  
1. Get token from header (Bearer Authorization).  
2. Check the token and extract the user ID.  
3. Generate a fingerprint from the request and hash it.  
4. Find this combination in Redis.  
5. Compare the token stored in Redis with the incoming token.  

Cool! When we register or log in, we create a record with the token.  

---

### Docker  

I packed the full application in a Docker Compose setup. I made an `internal-network` for Redis and PostgreSQL, and an `external-network` for services. Check the `.env` file and `init.sql` for database initialization.**  

**Note:** I know about secrets like `JWT_SECRET` or passwords, but this is only **the** first and simple version of **the** application. **In the next version, we will upgrade it and make it more secure and complicated.**  

---

Till next time!
