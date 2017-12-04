Record all the jobs that need to be run for a set of schema changes. The best way to approach this may be to run every schema change in a job, even if the schema change does not need to do anything more than MaybeIncrementVersion.

After this step, they can be run with BeginJob, as all of the jobs have already been created in the transaction. The client waits on each job to finish.

Because the client can cancel the session at any time, there needs to be a way to kickstart either:
1) The first schema change job to run, if none have been run yet.
2) The next schema change job, if there are multiple and one is currently running.

How to handle case 1 is not yet determined. For case 2, using job dependancies (i.e. add a "dependant" field on jobs), once the currently running job terminates it will search for all dependant jobs on itself, and set them to a status where they can be picked up by the registry. This can be integrated directly into the registry for OnSuccess/OnFailure. Alternatively, this can be implemented specifically for shcema changes within the Resumer OnFailureOrCancel or OnSuccess, where it looks up specific schema changer jobs to set to a beginning state. The alternative still requires a modifcation to jobs to add a "pending" status, where the job cannot be started yet, but has been created.
